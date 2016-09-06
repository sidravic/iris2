package client

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/satori/go.uuid"
	"github.com/Sirupsen/logrus"
	simplLogger "github.com/supersid/iris2/logger"
	"os"
	"runtime"
	"github.com/supersid/iris2/constants"
	"github.com/supersid/iris2/request"
	"fmt"
	"time"
	"errors"
)

const (
	POLL_TIMEOUT_FREQUENCY = 200 * time.Millisecond
)
var logger *logrus.Logger

type Client struct {
	brokerUrl string
	Socket    *zmq.Socket
	ID        string
	DebugMode bool
	poller    *zmq.Poller
}

func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	env := os.Getenv("IRIS_ENV")

	if env == "" {
		env = constants.DEVELOPMENT_ENV
	}

	logger = simplLogger.Init(env, "")
}

func (client *Client) setIdentity() {
	client.ID = uuid.NewV4().String()
}

func (client *Client) ReceiveMessage(waitFor time.Duration) (ClientResponse, error){
	sockets, err := client.poller.Poll(waitFor)

	if err != nil {
		logger.Error(fmt.Sprint("[client.go] Error while receiving messages: %s", err.Error()))
		return ClientResponse{}, err
	}

	if len(sockets) == 0 {
		err = errors.New("No Data")
		return ClientResponse{}, err
	}

	m, err := client.Socket.RecvMessage(0)

	if err != nil {
		logger.Error(fmt.Sprintf("[client.go] Error while receiving message from socket: %s", err.Error()))
		return ClientResponse{}, err
	}

	fmt.Println(m)
	res := NewClientResponse(m)

	return res, nil
}

func (client *Client) SendMessage(serviceName, message string){
	if message == "" || serviceName == ""{
		logger.Info(fmt.Sprintf("[client.go] Ignoring empty message sent by client ID: %s", client.ID))
		return
	}

	jsonMessage, err := request.CreateMessage(client.ID, constants.CLIENT_REQUEST, message,"", serviceName, "")
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"ServiceName": serviceName,
			"Message":     message,
			"ClientId":    client.ID,
		}).Error("[client.go] Could not marshal message.")
	}

	_, err = client.Socket.SendMessage(jsonMessage)

	if err != nil {
		logger.Error(fmt.Sprintf("[error.go] Error while sending message %s", err.Error()))
	}
}

func NewClient(brokerUrl string, sendWaterMark int) (*Client, error) {
	client := &Client{
		brokerUrl: brokerUrl,
		DebugMode: false,
	}

	socket, err := zmq.NewSocket(zmq.DEALER)

	if err != nil {
		logger.Error("[ERROR]: Could not create a new client socket.")
		return nil, err
	}

	client.Socket = socket
	client.Socket.SetSndhwm(sendWaterMark)
	client.setIdentity()

	client.poller = zmq.NewPoller()
	client.poller.Add(client.Socket, zmq.POLLIN)

	return client, err
}


func (client *Client) EnableDebugMode() bool{
	client.DebugMode = true
	return client.DebugMode
}



func (client *Client) Close() {
	client.Socket.Close()
}

func Start(brokerUrl string, sendWatermark int) (*Client){
	client, err := NewClient(brokerUrl, sendWatermark)

	err = client.Socket.Connect(brokerUrl)

	if err != nil {
		logger.Error(fmt.Sprintf("Client could not connect to broker. %s", err.Error()))
		panic(err)
	}

	logger.Info(fmt.Sprintf("Starting client id %s by connecting to %s", client.ID, client.brokerUrl))
	return client
}
