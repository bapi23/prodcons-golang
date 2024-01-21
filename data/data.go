package data

import (
	"fmt"
	"time"
)

func GenerateData() []byte {
	time.Sleep(time.Millisecond * 100)
	return []byte("generating data")
}

func ConsumeData(data []byte) {
	fmt.Println(string(data))
}
