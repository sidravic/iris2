package broker

import (

	"github.com/supersid/iris2/request"
	"github.com/supersid/iris2/service"
	"fmt"
	"errors"
)

type WorkerDisconnectHandler struct {
	errors  []error
	service  *service.Service
	broker   *Broker
	serviceWorker *service.ServiceWorker
}

func (broker *Broker) WorkerDisconnectHandler(req request.Request){
	handler := &WorkerDisconnectHandler{
		broker: broker,
		errors: make([]error, 0),
	}
	handler.findService(req.ServiceName)
	handler.findServiceWorker(req.ID)
	handler.deleteServiceWorker()

}

func (handler *WorkerDisconnectHandler) findService(serviceName string){
	err, service := handler.broker.FindService(serviceName)

	if err != nil {
		logger.Info(fmt.Sprintf("[worker_disconnect_handler.go] No service present with the name: %s", serviceName))
		handler.errors = append(handler.errors, errors.New("Service does not exist."))
	}

	handler.service = service
}

func (handler *WorkerDisconnectHandler) findServiceWorker(serviceWorkerID string){
	if len(handler.errors) == 0 {
		sw := handler.service.FindServiceWorker(serviceWorkerID)
		if sw.Identity != "" {
			handler.serviceWorker = sw
		}else{
			handler.errors = append(handler.errors, errors.New("Service Worker does not exist."))
		}

	}

}

func (handler *WorkerDisconnectHandler) deleteServiceWorker(){
	if len(handler.errors) == 0 {
		handler.service.DeleteWorker(handler.serviceWorker)
	}

}

