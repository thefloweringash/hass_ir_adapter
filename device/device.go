package device

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type Factory interface {
	New(c mqtt.Client, logger *log.Logger, stateDir string, emitter emitters.Emitter) (Device, error)
	Id() string
	EmitterId() string
}

type Device interface {
	Run() (func(), error)
}

type DeviceImpl interface {
	Config() map[string]interface{}
	DefaultState() State
	PushState(emitter emitters.Emitter, state State) error
}

type device struct {
	name    string
	state   state
	mqtt    Mqtt
	emitter emitters.Emitter
	impl    DeviceImpl
	logger  *log.Logger
}

type Update struct {
	binding Binding
	value   string
}

func New(
	c mqtt.Client,
	id string,
	name string,
	class string,
	emitter emitters.Emitter,
	impl DeviceImpl,
	stateDir string,
) (Device, error) {
	return &device{
		name: name,
		state: state{
			StateFile: path.Join(stateDir, id),
			State:     impl.DefaultState(),
		},
		mqtt: Mqtt{
			Class:  class,
			Id:     id,
			Client: c,
			Logger: log.New(os.Stdout, fmt.Sprintf("device/%s/mqtt: ", id), log.Lshortfile),
		},
		emitter: emitter,
		impl:    impl,
		logger:  log.New(os.Stdout, fmt.Sprintf("device/%s: ", id), log.Lshortfile),
	}, nil
}

func (d *device) publishConfig() error {
	config := d.impl.Config()
	return d.mqtt.PublishConfig(d.name, d.state.State, config)
}

func (d *device) Run() (func(), error) {
	stopChan := make(chan bool, 1)
	stopped := make(chan bool, 1)

	commandChan := make(chan Update)

	if err := d.mqtt.Subscribe(commandChan, d.impl.DefaultState()); err != nil {
		return nil, err
	}

	// Gather commands until the minimum interval of 100ms has elapsed
	readInitialUpdate := func() ([]Update, bool) {
		var updates []Update
		timeout := make(chan bool, 1)
		timeoutStarted := false

		for {
			select {
			case update := <-commandChan:
				updates = append(updates, update)
				if !timeoutStarted {
					timeoutStarted = true
					go func() {
						time.Sleep(100 * time.Millisecond)
						timeout <- true
					}()
				}

			case <-timeout:
				return updates, true

			case <-stopChan:
				return nil, false
			}
		}
	}

	// Continue to read updates until the channel is read
	readUpdatesUntil := func(ch chan bool) ([]Update, bool) {
		var updates []Update

		for {
			select {
			case update := <-commandChan:
				updates = append(updates, update)
			case <-ch:
				return updates, true
			case <-stopChan:
				return nil, false
			}
		}
	}

	if err := d.state.LoadState(); err != nil {
		d.logger.Printf("failed loading persisted state, using default")
	}

	if err := d.publishConfig(); err != nil {
		return nil, err
	}

	go func() {
		defer func() { stopped <- true }()

		defer func() {
			d.logger.Printf("removing config")
			if err := d.mqtt.RemoveConfig(); err != nil {
				d.logger.Printf("error removing retained config: %s", err)
			}
			d.logger.Printf("removing state")
			if err := d.mqtt.RemoveState(); err != nil {
				d.logger.Printf("error removing retained state: %s", err)
			}
		}()

		for {
			d.logger.Printf("state=%v, waiting for command", d.state.State)

			if err := d.mqtt.PublishState(d.state.State); err != nil {
				d.logger.Printf("error publishing state: %s", err)
			}

			updates, running := readInitialUpdate()
			if !running {
				return
			}
			d.logger.Printf("locking emitter for %d updates\n", len(updates))

			emitterReady := make(chan bool)
			go func() {
				d.emitter.Lock()
				emitterReady <- true
			}()

			moreUpdates, running := readUpdatesUntil(emitterReady)
			if !running {
				return
			}
			d.logger.Printf("got %d more updates while waiting for emitter", len(moreUpdates))

			updates = append(updates, moreUpdates...)

			newState := d.state.State
			for _, update := range updates {
				d.logger.Printf("processing update %s <- %s", update.binding, update.value)

				updatedState, err := update.binding.Apply(newState, update.value)
				if err != nil {
					d.logger.Printf("error deriving new state: %s", err)
					continue
				}
				newState = updatedState
			}

			d.logger.Printf("pushing state %v", newState)

			err := d.impl.PushState(d.emitter, newState)

			d.emitter.Unlock()

			if err != nil {
				d.logger.Printf("error %s pushing state: %v", err, newState)
				continue
			}

			d.state.State = newState

			if err := d.state.WriteState(); err != nil {
				d.logger.Printf("error persisting state: %s", err)
			}
		}
	}()

	return func() {
		stopChan <- true
		<-stopped
	}, nil
}
