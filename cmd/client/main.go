package main

import (
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
)

func main() {

	tc := client.CreateClient()

	tc.Read("test.txt")

	tc.Write("test1.txt", "Test File 1")
	time.Sleep(time.Second * 5)
	tc.Write("test2.txt", "Test File 2")
	time.Sleep(time.Second * 5)
	tc.Write("test3.txt", "Test File 3")
	time.Sleep(time.Second * 5)
	tc.Write("test4.txt", "Test File 4")
	time.Sleep(time.Second * 5)

	tc.Read("test1.txt")
	time.Sleep(time.Second * 5)
	tc.Read("test3.txt")
	time.Sleep(time.Second * 5)

}
