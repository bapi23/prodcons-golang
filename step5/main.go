package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
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

	go func() { // consumers goroutine
		var wg sync.WaitGroup
		wg.Add(NumOfConsumers)
		for i := 1; i <= NumOfConsumers; i++ {
			go func() {
				for {
					select {
					case d, ok := <-dataChan:
						if !ok {
							fmt.Println("closing consumer")
							wg.Done()
							return
						}
						data.ConsumeData(d)
					}
				}
			}()
		}
		wg.Wait()
		close(consumerCloseChan)
	}()

	go func() {
		var wg sync.WaitGroup
		wg.Add(NumOfConsumers)
		for i := 1; i <= NumOfProducers; i++ {
			go func() {
				for {
					select {
					case <-closeChan:
						fmt.Println("closing producer")
						wg.Done()
						return
					default:
						dataChan <- data.GenerateData()
					}
				}
			}()
		}
		wg.Wait()
		close(dataChan)
	}()

	<-consumerCloseChan
}
