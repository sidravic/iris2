package broker

import (
	"github.com/supersid/iris2/service"
	"github.com/supersid/iris2/request"
	"github.com/supersid/iris2/constants"

	"fmt"
)

type ClientRequestHandler struct{
	errors                []error
	broker                *Broker
	serviceAlreadyPresent bool
	service              *service.Service
}

func (broker *Broker) ClientRequestHandler(req request.Request){
	handler :=  &ClientRequestHandler{
		errors: make([]error, 0),
		broker: broker,
	}

	handler.findOrCreateService(req.ServiceName)
	handler.addServiceToBroker()
	handler.addServiceWorkerToService(req)
	handler.addClientRequestToService(req)
	err, req, serviceWorker := handler.processRequests()

	if err != nil {
		logger.Info(fmt.Sprintf("[client_request_handler.go] %s", err.Error()))
	}else{
		logger.Info(fmt.Sprintf("[client_request_handler.go] Processing request %s with service worker %s", req, serviceWorker))
		handler.broker.processClientRequest(req, serviceWorker)
	}

}

func (handler *ClientRequestHandler) findOrCreateService(serviceName string) {
	alreadyPresent, service := handler.broker.FindOrCreateService(serviceName)
	handler.serviceAlreadyPresent = alreadyPresent
	handler.service = service
}

func (handler *ClientRequestHandler) addServiceToBroker(){
	if !(handler.serviceAlreadyPresent){
		handler.broker.AddService(handler.service.ServiceName, handler.service)
	}
}

func (handler *ClientRequestHandler) addServiceWorkerToService(req request.Request){
	if !(handler.serviceAlreadyPresent){
		serviceWorker := service.NewServiceWorker(req.ID, req.Sender)
		handler.service.AddServiceWorker(serviceWorker)
	}
}

func (handler *ClientRequestHandler) addClientRequestToService(req request.Request) {
	handler.service.AddClientRequest(req)
}

func (handler *ClientRequestHandler) processRequests() (error, request.Request, *service.ServiceWorker){
	err, req, sw :=  handler.service.ProcessRequests()
	return err, req, sw
}

func (broker *Broker) processClientRequest(req request.Request, sw *service.ServiceWorker){
	logger.Info(fmt.Sprintf("[client_request_handler.go] Processing client request"))
	SocketIdCacheStore(req.ID, req.Sender)
	msg, err := request.CreateMessage("", constants.CLIENT_REQUEST_TO_WORKER, req.Data, "", req.ServiceName, req.ID)

	if err != nil {
		request.LogRequest(logger, req).Error("[client_request_handler.go] Unable to create message")
		return
	}

	logger.Debug(fmt.Sprintf("service worker %s", sw))
	logger.Debug(fmt.Sprintf("Request %s", req))

	_, err = broker.Socket.SendMessage(sw.Sender, msg)
	fmt.Println(err)
}




