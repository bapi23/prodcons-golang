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

func startProducers(ctx context.Context, dataChan chan []byte, num int) {
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
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
	// close data channel when we know there will be no more data produced
	close(dataChan)
}

func startConsumers(ctx context.Context, dataChan chan []byte, num int) {
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case d, ok := <-dataChan:
					if !ok {
						return
					}
					data.ConsumeData(d)
				}
			}
		}()
	}
	wg.Wait()
}

func main() {
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	dataChan := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			cancel()
		}
	}()

	go startConsumers(ctx, dataChan, NumOfProducers)
	startProducers(ctx, dataChan, NumOfProducers)
}
