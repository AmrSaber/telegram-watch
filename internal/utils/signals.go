package utils

import (
	"os"
	"os/signal"
)

func HandleInterrupt(handler func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			handler()
		}
	}()
}
