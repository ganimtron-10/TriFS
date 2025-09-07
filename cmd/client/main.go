package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
)

func main() {
	timeInterval := time.Second * 5
	var wg sync.WaitGroup

	tc := client.CreateClient()

	tc.Read("test.txt")

	numOfWriteFiles := 10
	for i := 0; i < numOfWriteFiles; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			tc.Write(fmt.Sprintf("test%d.txt", i), fmt.Sprintf("Test File %d", index))
		}(i)
		time.Sleep(timeInterval)
	}

	wg.Wait()
	time.Sleep(timeInterval)

	numOfReadFiles := 3
	for i := 0; i < numOfReadFiles; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tc.Read(fmt.Sprintf("test%d.txt", rand.Intn(numOfWriteFiles)))
		}()
		time.Sleep(timeInterval)
	}

	wg.Wait()

}
