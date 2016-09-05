package main

import (
	"github.com/supersid/iris2/client"

	"time"
	"fmt"
)

func main() {
	c := client.Start("tcp://127.0.0.1:5555")
	fmt.Println(fmt.Sprintf("Client ID: %s", c.ID))
	seq := 0
	for {
		time.Sleep(1 * time.Second)
		msg := fmt.Sprintf("Hello World! %d", seq)
		c.SendMessage("echo", msg)

		m, e := c.ReceiveMessage(200 * time.Millisecond)
		if e != nil {
			continue
		}
		fmt.Println("----------------------------------------")
		fmt.Println(fmt.Sprintf("Respose is: %s", m.ResponseData))
		fmt.Println("----------------------------------------")
		seq++
	}
}


