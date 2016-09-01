package request

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"fmt"
)

var logger *logrus.Logger

type Request struct {
	Sender   string
	ID       string
	Data     string
	Response string
	Command  string
}

func CreateMessage(senderId, command, data, responseData string ) ([]byte, error) {
	msg   := make([]string, 5)
	msg[0] = ""
	msg[1] = command
	msg[2] = senderId
	msg[3] = data
	msg[4] = ""

	jsonPayload, err := json.Marshal(msg)

	return jsonPayload, err
}


func UnWrapMessage(msg []string) (Request, error){
	sender  := msg[0]
	jsonPayload := msg[1]
	payload := make([]string, 6)

	err := json.Unmarshal([]byte(jsonPayload), &payload)

	if err != nil {
		return Request{}, err
	}

	for i, m := range payload {
		fmt.Printf("%d. %s", i, m)
		fmt.Println("")
	}
	req := Request{
		Sender:  sender,
		Command: payload[1],
		ID:      payload[2],
		Data:    payload[3],
		Response:"",
	}
	return req, nil
}


