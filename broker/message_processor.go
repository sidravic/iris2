package broker

import (
	"github.com/supersid/iris2/request"
	"github.com/supersid/iris2/constants"
	"fmt"
	"github.com/supersid/iris2/service"
)

type ProcessMessage interface {
	ProcessMessage(request.Request)
}

func (broker *Broker) ProcessMessage(req request.Request) {
	if broker.DebugMode {
		logger.Debug(fmt.Sprintf("[message_processor.go] Processing message %s", req))
	}

	switch req.Command {
	case constants.WORKER_READY:
		broker.WorkerReadyHandler(req)
	case constants.CLIENT_REQUEST:
		broker.ClientRequestHandler(req)
	case constants.WORKER_RESPONSE:
		broker.WorkerResponseHandler(req)
	default:
		panic("OMFG")
	}
}

func (broker *Broker) FindOrCreateService(serviceName string) (bool, *service.Service) {
	alreadyPresent := false
	var s *service.Service

	for srvName, service := range broker.Services {
		if srvName == serviceName {
			alreadyPresent = true
			s = service
			break;
		}
	}

	if !alreadyPresent {
		s = &service.Service{
			ServiceName:    serviceName,
			ClientRequest:  make([]request.Request,0),
			WaitingWorkers: make([]*service.ServiceWorker, 0),
			Workers:        make(map[string]*service.ServiceWorker),
		}

		logger.WithFields(map[string]interface{}{
			"ServiceName": serviceName,
		}).Info("New service created.")
	}

	return alreadyPresent, s
}

func (broker *Broker) AddService (serviceName string, s *service.Service){
	broker.Services[serviceName] = s
}