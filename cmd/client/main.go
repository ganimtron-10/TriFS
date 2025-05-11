package main

import "github.com/ganimtron-10/TriFS/internal/client"

func main() {
	trifsClient := client.CreateClient()
	trifsClient.Read("test1.txt")
}
