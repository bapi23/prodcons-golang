package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bapi23/prodcons-golang/data"
)

func main() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	dataChan := make(chan []byte, 1)

	go func() {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			close(dataChan)
		}
	}()

	go func() { // consumer goroutine
		for d := range dataChan {
			data.ConsumeData(d)
		}
	}()

	for { // producer goroutine
		dataChan <- data.GenerateData()
	}
}
