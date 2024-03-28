package data

import (
	"fmt"
	"time"
)

func GenerateData() []byte {
	time.Sleep(time.Second)
	return []byte("generated data")
}

func ConsumeData(data []byte) {
	fmt.Println(string(data))
}
