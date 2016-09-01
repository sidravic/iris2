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
	HEARTBEAT_INTERVAL = 2 * time.Second
	POLL_FREQUENCY     = 100 * time.Millisecond
)

var logger *logrus.Logger
var env string

type Worker struct{
	ID 	        string
	BrokerUrl       string
	ServiceName     string
	LastHeartBeatAt time.Time
	socket          *zmq.Socket
	DebugMode       bool
}

func NewWorker(brokerUrl string, serviceName string, env string) (*Worker, error){
	worker := &Worker{ID:          	  uuid.NewV4().String(),
		          BrokerUrl:   	  brokerUrl,
		          ServiceName: 	  serviceName,
			  LastHeartBeatAt:time.Now(),
			  DebugMode:      false,
		}

	socket, err := zmq.NewSocket(zmq.DEALER)

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	worker.socket = socket

	if env == constants.DEVELOPMENT_ENV || env == constants.TEST_ENV {
		worker.DebugMode = true
	}

	return worker, nil
}

func (worker *Worker) SendReadyMessage() (bool){
	sendStatus := true

	command  := constants.WORKER_READY
	data     := worker.ServiceName
	SenderID := worker.ID
	response := ""

	msg, err := request.CreateMessage(SenderID, command, data, response)
	if err != nil {
		sendStatus = false
		logger.Error("Message creation error: %s", err.Error())

	}

	logger.Info(msg)
	_, err = worker.socket.SendMessage(msg)

	if err != nil {
		sendStatus = false
		logger.Error(fmt.Sprintf("Error while sending WORKER_READY message: %s", err.Error()))
	}

	logger.Debug("ReadyMessage Sent.")
	return sendStatus
}

func (worker *Worker) Process(messageChannel chan []string){
	poller := zmq.NewPoller()
	poller.Add(worker.socket, zmq.POLLIN)

	for {
		if worker.DebugMode {
			logger.Debug("Looping...")
		}
		sendStatus := worker.SendReadyMessage()

		if !sendStatus {
			continue
		}

		incomingSockets, err := poller.Poll(POLL_FREQUENCY)

		if err != nil{
			logger.Error(fmt.Sprintf("Error while polling: %s", err.Error()))
		}

		if len(incomingSockets) > 0 {
			msg, err := worker.socket.RecvMessage(0)

			if err != nil {
				logger.Error(fmt.Sprintf("Error while receving messages: %s", err.Error()))
			}

			messageChannel <- msg
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

func Start(brokerUrl string, serviceName string) (chan []string){
	worker, err := NewWorker(brokerUrl, serviceName, env)

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	err = worker.socket.Connect(brokerUrl)

	if err != nil {
		logger.Error(fmt.Sprintf("Bind Error: %s", err.Error()))
		panic(err)
	}

	messageChannel := make(chan []string)
	go worker.Process(messageChannel)
	return messageChannel

}
