package main

import (
	"github.com/supersid/iris2/worker"
	"fmt"
)

func main() {
	channel := worker.Start("tcp://127.0.0.1:5555", "echo")
	for  {
		msg := <-channel
		fmt.Println(msg)
	}

}
