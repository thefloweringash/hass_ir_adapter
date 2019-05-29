package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/config"
	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s", sig)
		done <- true
	}()

	<-done
}

func run(cfg config.Config, stateDir string) error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "amnesiac"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTT.Broker)
	opts.SetClientID(fmt.Sprintf("hass_ir_adapter-%d-%v", os.Getpid(), hostname))
	opts.SetUsername(cfg.MQTT.Username)
	opts.SetPassword(cfg.MQTT.Password)
	opts.SetCleanSession(false)

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	emitters := map[string]emitters.Emitter{}
	for _, emitterFactory := range cfg.Emitters {
		emitter, err := emitterFactory.New(c)
		if err != nil {
			return err
		}
		emitters[emitterFactory.Id()] = emitter
	}

	devices := map[string]device.Device{}
	for _, deviceFactory := range cfg.Devices {
		emitter := emitters[deviceFactory.EmitterId()]

		if emitter == nil {
			return fmt.Errorf("invalid emitter reference '%s', in device '%s'",
				deviceFactory.EmitterId(), deviceFactory.Id())
		}

		aircon, err := deviceFactory.New(c, emitter, stateDir)
		if err != nil {
			return err
		}

		devices[deviceFactory.Id()] = aircon
	}

	for _, device := range devices {
		stop, err := device.Run()
		if err != nil {
			return err
		}
		defer stop()
	}

	waitForSignal()

	return nil
}

func main() {
	// mqtt.DEBUG = log.New(os.Stdout, "", log.LstdFlags)
	mqtt.ERROR = log.New(os.Stdout, "", log.Lshortfile)

	var configFile, stateDir string
	flag.StringVar(&configFile, "config-file", "", "Configuration yaml path")
	flag.StringVar(&stateDir, "state-dir", "", "StatePtr directory")
	flag.Parse()

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	if err := run(*cfg, stateDir); err != nil {
		panic(err)
	}
}
