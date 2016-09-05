package client



import (
	"encoding/json"
)

type ClientResponse struct {
	RequestID             string
	OriginalClientSenderID string
	Command               string
	Data		      string
	ServiceName           string
	ResponseData          string
}

func NewClientResponse(msg []string) ClientResponse{
	m := make([]string, 6)

	json.Unmarshal([]byte(msg[0]), &m)

	cr  := ClientResponse{
		OriginalClientSenderID: m[0],
		Command:                m[1],
		Data:                   m[3],
		ResponseData:           m[4],
		ServiceName:            m[5],
		RequestID:              m[6],
	}


	return cr
}

