package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bapi23/prodcons-golang/data"
)

func main() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	dataChan := make(chan []byte, 1)
	closeChan := make(chan struct{})
	consumerCloseChan := make(chan struct{})

	go func() {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			close(closeChan)
		}
	}()

	go func() { // consumer goroutine
		for {
			select {
			case d, ok := <-dataChan:
				if !ok {
					fmt.Println("closing consumer")
					close(consumerCloseChan)
					return
				}
				data.ConsumeData(d)
			}
		}
	}()

	go func() { // producer goroutine
		for {
			select {
			case <-closeChan:
				fmt.Println("closing producer")
				close(dataChan)
				return
			default:
				dataChan <- data.GenerateData()
			}
		}
	}()

	<-consumerCloseChan
}
