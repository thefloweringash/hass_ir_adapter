package irkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type IRKit struct {
	Endpoint string
	Logger   *log.Logger
	Token    *token
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

func (irkit *IRKit) Emit(commands ...emitters.Command) error {
	token, delay := irkit.Token.Take()
	defer irkit.Token.Return(token)

	irkit.Logger.Printf("Acquired token after %s delay", delay)

	var intervals []uint16
	for _, cmd := range commands {
		intervals = append(intervals, cmd.Intervals()...)
	}

	for i := range intervals {
		intervals[i] *= 2
	}

	msg := Message{
		Format:    "raw",
		Freq:      38,
		Intervals: intervals,
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

	irkit.Logger.Printf("Sending request: %s\n", jsonString)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Unexpected status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	irkit.Logger.Printf("Message response: %s\n", string(body))

	return nil
}
