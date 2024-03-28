package main

import (
	"github.com/bapi23/prodcons-golang/data"
)

func main() {
	dataChan := make(chan []byte)

	go func() { // consumer goroutine
		for d := range dataChan {
			data.ConsumeData(d)
		}
	}()

	for { // producer goroutine
		dataChan <- data.GenerateData()
	}
}
