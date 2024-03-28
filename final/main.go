package main

import (
	"context"
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
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	dataChan := make(chan []byte, 1)
	consumersClosedChan := make(chan struct{})

	go runConsumers(dataChan, consumersClosedChan)
	go runProducers(ctx, dataChan)

	<-consumersClosedChan
}

func runConsumers(dataChan chan []byte, consumersClosedChan chan struct{}) {
	var wg sync.WaitGroup
	wg.Add(NumOfConsumers)
	for i := 0; i < NumOfConsumers; i++ {
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
	close(consumersClosedChan)
}

func runProducers(ctx context.Context, dataChan chan []byte) {
	var wg sync.WaitGroup
	wg.Add(NumOfConsumers)
	for i := 0; i < NumOfProducers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
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
}
