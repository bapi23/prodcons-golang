for i := 0; i < NumOfConsumers; i++ {
	go func() { // consumer goroutine
		for {
			select {
			case d := <-dataChan:
				data.ConsumeData(d)
			case <-closeChan:
				for range dataChan {
				}
				fmt.Println("closing consumer")
				return
			}
		}
	}()
}