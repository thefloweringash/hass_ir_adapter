package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/aircon/mitsubishi_gp82"
	"github.com/thefloweringash/hass_ir_adapter/config"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\nMSG:%s\n", msg.Topic(), msg.Payload())
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		sig := <-sigs
		fmt.Printf("Received signal: %s", sig)
		done <- true
	}()

	<-done
}

func makeAirconController(airconType string, emitter emitters.Emitter) (aircon.AirconController, error) {
	var impl aircon.AirconController
	switch airconType {
	case "mitsubishi_gp82":
		impl = mitsubishi_gp82.NewAircon(emitter)
	default:
		return nil, fmt.Errorf("Unknown aircon type: %s", airconType)
	}
	return impl, nil
}

func makeEmitter(client mqtt.Client, sendTopic string) (emitters.Emitter, error) {
	return irblaster.NewMQTTIRBlaster(client, sendTopic), nil
}

func run(cfg config.Config, stateDir string) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTT.Broker)
	opts.SetDefaultPublishHandler(f)
	opts.SetClientID("hass_ir_adapter")
	opts.SetUsername(cfg.MQTT.Username)
	opts.SetPassword(cfg.MQTT.Password)
	opts.SetCleanSession(false)

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	emitters := map[string]emitters.Emitter{}
	for _, emitterCfg := range cfg.Emitters {
		emitter, err := makeEmitter(c, emitterCfg.Topic)
		if err != nil {
			return err
		}
		emitters[emitterCfg.Id] = emitter
	}

	aircons := map[string]*aircon.Aircon{}
	for _, airconCfg := range cfg.Aircons {
		emitter := emitters[airconCfg.Emitter]

		if emitter == nil {
			return fmt.Errorf("Invalid emitter reference '%s', in aircon '%s'",
				airconCfg.Emitter, airconCfg.Id)
		}

		impl, err := makeAirconController(airconCfg.Type, emitter)
		if err != nil {
			return err
		}

		aircon, err := aircon.NewAircon(
			c, emitter, impl,
			airconCfg.Id, airconCfg.Name, airconCfg.TemperatureTopic,
			stateDir,
		)
		if err != nil {
			return err
		}

		aircons[airconCfg.Id] = aircon
	}

	for _, aircon := range aircons {
		stop, err := aircon.Run()
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
	mqtt.ERROR = log.New(os.Stdout, "", log.LstdFlags)

	var configFile, stateDir string
	flag.StringVar(&configFile, "config-file", "", "Configuration yaml path")
	flag.StringVar(&stateDir, "state-dir", "", "State directory")
	flag.Parse()

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	if err := run(*cfg, stateDir); err != nil {
		panic(err)
	}
}
