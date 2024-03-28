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
}

func startConsumers(ctx context.Context, dataChan chan []byte, num int) {
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case d := <-dataChan:
					data.ConsumeData(d)
				case <-ctx.Done():
					select {
					case d := <-dataChan:
						data.ConsumeData(d)
					default:
						return
					}
					fmt.Println("closing consumer")
					wg.Done()
					return
				}
			}
		}()
	}
	wg.Wait()
	return
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	dataChan := make(chan []byte)

	go startConsumers(ctx, dataChan, NumOfProducers)
	startProducers(ctx, dataChan, NumOfProducers)
}
