package request

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"fmt"
	"github.com/satori/go.uuid"
)

var logger *logrus.Logger

type Request struct {
	RequestId   string
	Sender      string
	ID          string
	Data        string
	Response    string
	Command     string
	ServiceName string
}

func CreateMessage(senderId, command, data, responseData, serviceName, originalSenderId string) ([]byte, error) {
	msg   := make([]string, 7)
	msg[0] = originalSenderId
	msg[1] = command
	msg[2] = senderId
	msg[3] = data
	msg[4] = responseData
	msg[5] = serviceName
	msg[6] = fmt.Sprintf("RID-%s", uuid.NewV4().String())

	jsonPayload, err := json.Marshal(msg)

	return jsonPayload, err
}


func UnWrapMessage(msg []string) (Request, error){
	sender  := msg[0]
	jsonPayload := msg[1]
	fmt.Println(sender)
	fmt.Println(jsonPayload)

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
		RequestId:   payload[6],
		Sender:      sender,
		Command:     payload[1],
		ID:          payload[2],
		Data:        payload[3],
		Response:    payload[4],
		ServiceName: payload[5],
	}
	return req, nil
}


