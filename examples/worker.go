package main

import (
	"github.com/supersid/iris2/worker"
	"fmt"


)

func main() {
	channel, w := worker.Start("tcp://127.0.0.1:5555", "echo")
	for  {
		wr := <-channel
		fmt.Println(wr)
		response := fmt.Sprintf("%s-Response",wr.Data)
		wr.ResponseData = response
		w.SendResponse(wr)
	}

}
