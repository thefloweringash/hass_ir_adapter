package device

import (
	"fmt"
	"log"
	"os"
	"path"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type Factory interface {
	New(c mqtt.Client, emitter emitters.Emitter, stateDir string) (Device, error)
	Id() string
	EmitterId() string
}

type Device interface {
	Run() (func(), error)
}

type DeviceImpl interface {
	Config() interface{}
	DefaultState() State
	PushState(state State) error
}

type device struct {
	name   string
	state  state
	mqtt   Mqtt
	impl   DeviceImpl
	logger *log.Logger
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
		impl:   impl,
		logger: log.New(os.Stdout, fmt.Sprintf("device/%s: ", id), log.Lshortfile),
	}, nil
}

func (d *device) publishConfig() error {
	config := d.impl.Config()
	return d.mqtt.PublishConfig(d.name, d.state.State, config)
}

func (d *device) Run() (func(), error) {
	stopChan := make(chan bool, 1)
	stopDoneChan := make(chan bool, 1)

	commandChan := make(chan Update)

	if err := d.mqtt.Subscribe(commandChan, d.impl.DefaultState()); err != nil {
		return nil, err
	}

	if err := d.state.LoadState(); err != nil {
		d.logger.Printf("failed loading persisted state, using default")
	}

	if err := d.publishConfig(); err != nil {
		return nil, err
	}

	go func() {
		running := true
		for running == true {
			d.logger.Printf("state=%v, waiting for command", d.state.State)

			if err := d.mqtt.PublishState(d.state.State); err != nil {
				d.logger.Printf("Error publishing state: %s", err)
			}

			var update Update
			select {
			case update = <-commandChan:
			case <-stopChan:
				running = false
				continue
			}

			d.logger.Printf("processing update %s <- %s", update.binding, update.value)

			newState, err := update.binding.Apply(d.state.State, update.value)
			if err != nil {
				d.logger.Printf("Error deriving new state: %s", err)
				continue
			}

			if err := d.impl.PushState(newState); err != nil {
				d.logger.Printf("Error %s pushing state: %v", err, newState)
				continue
			}

			d.state.State = newState

			if err := d.state.WriteState(); err != nil {
				d.logger.Printf("Error persisting state: %s", err)
			}
		}

		if err := d.mqtt.RemoveConfig(); err != nil {
			d.logger.Printf("Error removing retained config: %s", err)
		}
		if err := d.mqtt.RemoveState(); err != nil {
			d.logger.Printf("Error remove retained state: %s", err)
		}

		stopDoneChan <- true
	}()

	return func() {
		stopChan <- true
		<-stopDoneChan
	}, nil
}
