package main

import (
	"fmt"
	"time"

	"github.com/ganimtron-10/TriFS/internal/client"
	"github.com/ganimtron-10/TriFS/internal/master"
)

func main() {
	go master.Init()

	fmt.Println("Waiting 5 sec to get Master up")
	time.Sleep(time.Second * 5)

	tc := client.CreateClient()
	tc.Read("test.txt")
}
