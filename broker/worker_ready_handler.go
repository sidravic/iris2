package broker

import (
	"github.com/supersid/iris2/request"
	"github.com/supersid/iris2/service"
)



type WorkerReadyHandler struct {
	errors                []error
	service               *service.Service
	serviceAlreadyPresent bool
	broker                *Broker
}

type WorkerReadyHandlersIntF interface  {
	WorkerReadyHandler(request.Request)
}

func (broker *Broker) WorkerReadyHandler(req request.Request) {
	handler := &WorkerReadyHandler{
		errors: make([]error, 1),
		broker: broker,
	}
	handler.findOrCreateService(req.Data)
	handler.AddServiceToBroker()
	handler.AddServiceWorkerToService(req)
}

func (handler *WorkerReadyHandler) findOrCreateService(serviceName string) {
	broker := handler.broker
	alreadyPresent, service := broker.FindOrCreateService(serviceName)
	handler.service = service
	handler.serviceAlreadyPresent = alreadyPresent

}

func (handler *WorkerReadyHandler) AddServiceToBroker(){
	if !(handler.serviceAlreadyPresent) {
		handler.broker.AddService(handler.service.ServiceName, handler.service)
	}
}


func (handler *WorkerReadyHandler) AddServiceWorkerToService(req request.Request){
	if !(handler.serviceAlreadyPresent) {
		serviceWorker := service.NewServiceWorker(req.ID, req.Sender)
		handler.service.AddServiceWorker(serviceWorker)
	}

}