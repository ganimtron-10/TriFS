package main

import (
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
)

func main() {
	timeInterval := time.Second * 5

	tc := client.CreateClient()

	tc.Read("test.txt")

	tc.Write("test1.txt", "Test File 1")
	time.Sleep(timeInterval)
	tc.Write("test2.txt", "Test File 2")
	time.Sleep(timeInterval)
	tc.Write("test3.txt", "Test File 3")
	time.Sleep(timeInterval)
	tc.Write("test4.txt", "Test File 4")
	time.Sleep(timeInterval)

	tc.Read("test1.txt")
	time.Sleep(timeInterval)
	tc.Read("test3.txt")
	time.Sleep(timeInterval)

}
