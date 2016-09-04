package main

import (
	"github.com/supersid/iris2/client"

	"time"
	"fmt"
)

func main() {
	c := client.Start("tcp://127.0.0.1:5555")

	seq := 0
	for {
		time.Sleep(3 * time.Second)
		msg := fmt.Sprintf("Hello World! %d", seq)
		c.SendMessage("echo", msg)
		fmt.Println(msg)
		seq++
	}
}


