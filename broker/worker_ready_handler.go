package broker

import (
	"github.com/supersid/iris2/request"
	"github.com/supersid/iris2/service"
	"fmt"
)



type WorkerReadyHandler struct {
	errors                []error
	service               *service.Service
	serviceAlreadyPresent bool
	broker                *Broker
	serviceWorker         *service.ServiceWorker
}

type WorkerReadyHandlersIntF interface  {
	WorkerReadyHandler(request.Request)
}

func (broker *Broker) WorkerReadyHandler(req request.Request) {
	handler := &WorkerReadyHandler{
		errors: make([]error, 1),
		broker: broker,
	}
	handler.findOrCreateService(req.ServiceName)
	handler.addServiceToBroker()
	handler.addServiceWorkerToService(req)
	handler.updateServiceWorkerHeartBeatExpiry()
	err, req, serviceWorker := handler.processRequests()

	if err != nil {
		logger.Info(fmt.Sprintf("[client_request_handler.go] %s", err.Error()))
	}else{
		logger.Info(fmt.Sprintf("[client_request_handler.go] Processing request %s with service worker %s", req, serviceWorker))
		handler.broker.processClientRequest(req, serviceWorker)
	}
}

func (handler *WorkerReadyHandler) findOrCreateService(serviceName string) {
	broker := handler.broker
	alreadyPresent, service := broker.FindOrCreateService(serviceName)
	handler.service = service
	handler.serviceAlreadyPresent = alreadyPresent

}

func (handler *WorkerReadyHandler) addServiceToBroker(){
	if !(handler.serviceAlreadyPresent) {
		handler.broker.AddService(handler.service.ServiceName, handler.service)
	}
}


func (handler *WorkerReadyHandler) addServiceWorkerToService(req request.Request){
	serviceWorker := service.NewServiceWorker(req.ID, req.Sender)
	handler.serviceWorker = serviceWorker
	handler.service.AddServiceWorker(serviceWorker)
}

func (handler *WorkerReadyHandler) updateServiceWorkerHeartBeatExpiry(){
	handler.serviceWorker.UpdateExpiry()
}

func (handler *WorkerReadyHandler) processRequests() (error, request.Request, *service.ServiceWorker){
	err, req, sw :=  handler.service.ProcessRequests()
	return err, req, sw
}