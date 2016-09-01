package main

import "github.com/supersid/iris2/broker"

func main(){
	broker.Start("tcp://*:5555")
}
