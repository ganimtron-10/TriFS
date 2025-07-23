package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
)

func main() {
	timeInterval := time.Second * 5

	tc := client.CreateClient()

	tc.Read("test.txt")

	numOfWriteFiles := 10
	for i := 0; i < numOfWriteFiles; i++ {
		go tc.Write(fmt.Sprintf("test%d.txt", i), fmt.Sprintf("Test File %d", i))
		time.Sleep(timeInterval)
	}

	time.Sleep(timeInterval)

	numOfReadFiles := 3
	for i := 0; i < numOfReadFiles; i++ {
		go tc.Read(fmt.Sprintf("test%d.txt", rand.Intn(numOfWriteFiles)))
		time.Sleep(timeInterval)
	}

}
