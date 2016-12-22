package main

import (
	"fmt"
	"os"
	"terrigenesis/fileclient/client"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: username")
	} else {
		client.StartClient(os.Args[1:])
	}
}
