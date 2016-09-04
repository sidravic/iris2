package worker

import (
	"fmt"
	"encoding/json"
)

type WorkerRequest struct {
	RequestID             string
	OriginalClientSenderID string
	Command               string
	Data		      string
        ServiceName           string
        ResponseData          string
}

func NewWorkerRequest(msg []string) WorkerRequest{
	fmt.Println(msg[0])
	m := make([]string, 6)
	json.Unmarshal([]byte(msg[0]), &m)

	wr  := WorkerRequest{
		OriginalClientSenderID: m[0],
		Command:                m[1],
		Data:                   m[3],
		ResponseData:           "",
		ServiceName:            m[5],
		RequestID:              m[6],
	}


	return wr
}
