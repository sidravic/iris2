package worker

import (
	simplLogger "github.com/supersid/iris2/logger"
	"github.com/supersid/iris2/constants"
	"github.com/supersid/iris2/request"
	uuid "github.com/satori/go.uuid"
	zmq "github.com/pebbe/zmq4"
	"github.com/Sirupsen/logrus"
	"runtime"
	"time"
	"fmt"
	"os"
)

const (
	HEARTBEAT_INTERVAL = 3 * time.Second * 3
	POLL_FREQUENCY     = 100 * time.Millisecond
)

var logger *logrus.Logger
var env string
var brokerUrl string
var serviceName string

type Worker struct{
	ID 	        string
	BrokerUrl       string
	ServiceName     string
	HeartBeatExpiry	time.Time
	Socket          *zmq.Socket
	DebugMode       bool
}

func NewWorker(brokerUrl, serviceName, env string) (*Worker, error){
	worker := &Worker{ID:          	  uuid.NewV4().String(),
		          BrokerUrl:   	  brokerUrl,
		          ServiceName: 	  serviceName,
			  HeartBeatExpiry:time.Now().Add(HEARTBEAT_INTERVAL),
			  DebugMode:      false,
		}

	socket, err := zmq.NewSocket(zmq.DEALER)

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	worker.Socket = socket
	worker.Socket.SetRcvhwm(constants.WORKER_RECV_HWM)

	if env == constants.DEVELOPMENT_ENV || env == constants.TEST_ENV {
		worker.DebugMode = true
	}

	return worker, nil
}

func Connect(brokerUrl, serviceName string) (*Worker){
	worker, err := NewWorker(brokerUrl, serviceName, env)

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	err = worker.Socket.Connect(brokerUrl)

	if err != nil {
		logger.Error(fmt.Sprintf("Bind Error: %s", err.Error()))
		panic(err)
	}

	return worker
}

func (worker *Worker) SendReadyMessage() (bool){
	sendStatus := true

	command     := constants.WORKER_READY
	data        := worker.ServiceName
	SenderID    := worker.ID
	response    := ""
	serviceName := worker.ServiceName

	msg, err := request.CreateMessage(SenderID, command, data, response, serviceName, "")

	if err != nil {
		sendStatus = false
		logger.Error("[worker.go] Message creation error: %s", err.Error())
		return sendStatus
	}

	_, err = worker.Socket.SendMessage(msg)

	if err != nil {
		sendStatus = false
		logger.Error(fmt.Sprintf("[worker.go] Error while sending WORKER_READY message: %s", err.Error()))
	}

	return sendStatus
}

func (worker *Worker) SendResponse(wr WorkerRequest) (bool){
	var sendStatus bool
	command          := constants.WORKER_RESPONSE
	data             := wr.Data
	senderID         := worker.ID
	response         := wr.ResponseData
	originalSenderId := wr.OriginalClientSenderID
	serviceName      := wr.ServiceName

	msg, err := request.CreateMessage(senderID, command, data, response, serviceName, originalSenderId)

	if err != nil {
		sendStatus = false
		logger.Error("[worker.go] Message creation error: %s", err.Error())
		return sendStatus
	}

	_, err = worker.Socket.SendMessage(msg)

	if err != nil {
		sendStatus = false
		logger.WithFields(map[string]interface{}{
			"Command":          command,
			"data":             data,
			"senderID":         senderID,
			"response":         response,
			"originalSenderId": originalSenderId,
			"serviceName":      serviceName,
		}).Error(fmt.Sprintf("[worker.go] Error while sending WORKER_RESPONSE message: %s", err.Error()))
	}

	return true
}


func (worker *Worker) SendDisconnect() bool{
	sendStatus := true

	command     := constants.WORKER_DISCONNECT
	data        := worker.ServiceName
	SenderID    := worker.ID
	response    := ""
	serviceName := worker.ServiceName

	msg, err := request.CreateMessage(SenderID, command, data, response, serviceName, "")

	if err != nil {
		sendStatus = false
		logger.Error("[worker.go] Message creation error: %s", err.Error())
		return sendStatus
	}

	_, err = worker.Socket.SendMessage(msg)

	if err != nil {
		sendStatus = false
		logger.Error(fmt.Sprintf("[worker.go] Error while sending WORKER_DISCONNECT message: %s", err.Error()))
	}

	return sendStatus
}

func (worker *Worker) Process(messageChannel chan WorkerRequest){
LOOP:
	poller := zmq.NewPoller()
	poller.Add(worker.Socket, zmq.POLLIN)

	for {
		sendStatus := worker.SendReadyMessage()

		if !sendStatus {
			continue
		}

		incomingSockets, err := poller.Poll(POLL_FREQUENCY)

		if err != nil{
			logger.Error(fmt.Sprintf("Error while polling: %s", err.Error()))
		}

		if len(incomingSockets) > 0 {
			msg, err := worker.Socket.RecvMessage(0)

			if err != nil {
				logger.Error(fmt.Sprintf("Error while receving messages: %s", err.Error()))
			}

			if len(msg) > 0 {
				worker.ProcessMessage(msg, messageChannel)

			}

		} else {
			isReconnect := worker.IsDisconnectAndReconnect()
			if isReconnect {
				logger.Info("Reconnection with broker. Broker appears to be unreachable.")
				sendStatus := worker.SendDisconnect()
				if !sendStatus {
					panic("Could not disconnect")
				}
				w := worker
				worker = nil
				logger.WithFields(map[string]interface{}{
					"BrokerUrl": w.BrokerUrl,
					"ServiceName":w.ServiceName,
				})
				worker = Connect(w.BrokerUrl, w.ServiceName)
				goto LOOP
			}
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	env = os.Getenv("IRIS_ENV")

	if env == "" {
		env = constants.DEVELOPMENT_ENV
	}

	logger = simplLogger.Init(env, "")
}

func GetLogger() *logrus.Logger{
	return logger;
}

func Start(brokerUrl string, serviceName string) (chan WorkerRequest, *Worker){
	worker := Connect(brokerUrl, serviceName)
	messageChannel := make(chan WorkerRequest)
	go worker.Process(messageChannel)
	return messageChannel, worker

}
