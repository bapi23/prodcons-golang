package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bapi23/prodcons-golang/data"
)

const NumOfProducers = 3
const NumOfConsumers = 3

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

	for i := 0; i < NumOfConsumers; i++ {
		go func() {
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
	}

	for i := 0; i < NumOfProducers; i++ {
		go func() {
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
	}

	<-consumerCloseChan
}
