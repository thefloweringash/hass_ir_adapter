package irkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/encodings"
)

type IRKit struct {
	Endpoint string
	Logger   *log.Logger
	Token    *token
}

func (irkit *IRKit) Lock() {
	delay := irkit.Token.Take()
	irkit.Logger.Printf("acquired token after %s delay", delay)
}

func (irkit *IRKit) Unlock() {
	irkit.Token.Return()
}

func New(endpoint string, logger *log.Logger) *IRKit {
	// Enforce a 2 second delay between requests
	// https://github.com/irkit/device/issues/10
	token := NewToken(2 * time.Second)

	return &IRKit{
		Endpoint: endpoint,
		Logger:   logger,
		Token:    token,
	}
}

type Message struct {
	Format    string   `json:"format"`
	Freq      int      `json:"freq"`
	Intervals []uint16 `json:"data"`
}

func (irkit *IRKit) Emit(command emitters.Command) error {
	intervalCommand, ok := command.(encodings.ToIntervals)
	if !ok {
		return errors.New("command not convertible to intervals")
	}

	intervals := intervalCommand.ToIntervals()
	irkitIntervals := []uint16{}

	irkit.Logger.Printf("raw intervals: %v\n", intervals)

	for i := range intervals {
		x := int(intervals[i]) * 2
		for x > math.MaxUint16 {
			irkitIntervals = append(irkitIntervals, math.MaxUint16, 0)
			x -= math.MaxUint16
		}
		irkitIntervals = append(irkitIntervals, uint16(x))
	}

	msg := Message{
		Format:    "raw",
		Freq:      38,
		Intervals: irkitIntervals,
	}

	var err error

	jsonString, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", irkit.Endpoint+"/messages", bytes.NewBuffer(jsonString))

	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("X-Requested-With", "curl")

	irkit.Logger.Printf("sending request: %s\n", jsonString)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	irkit.Logger.Printf("message response: %s\n", string(body))

	return nil
}
