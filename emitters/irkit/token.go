package irkit

import (
	"time"
)

// This token is something that allows consumers to enforce a minimum
// interval between actions. When a token is Return()ed, it will not
// be available to Take()en again until minimumGap has passed.

type token struct {
	minimumGap   time.Duration
	tokenChannel chan int
}

func NewToken(minimumGap time.Duration) *token {
	tokenChannel := make(chan int, 1)
	tokenChannel <- 0
	return &token{
		minimumGap:   minimumGap,
		tokenChannel: tokenChannel,
	}
}

func (token *token) Take() (int, time.Duration) {
	var value int
	select {
	case value = <-token.tokenChannel:
		return value, 0
	default:
		waitStart := time.Now()
		value = <-token.tokenChannel
		return value, time.Now().Sub(waitStart)
	}
}

func (token *token) Return(val int) {
	go func() {
		time.Sleep(token.minimumGap)
		token.tokenChannel <- val
	}()
}
