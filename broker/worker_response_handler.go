package broker

import (
	"github.com/supersid/iris2/request"
	"errors"
	"github.com/supersid/iris2/constants"
	"fmt"
)

type WorkerResponseHandler struct {
	errors       []error
	broker       *Broker
	clientSender string
}

func (broker *Broker) WorkerResponseHandler(req request.Request){
	handler := &WorkerResponseHandler{
		broker: broker,
		errors: make([]error, 0),
	}

	handler.findClientSender(req)
	handler.relayResponse(req)
}

func (handler *WorkerResponseHandler) findClientSender(req request.Request){
	logger.Info("----------------------------")
	logger.Info(req)
	logger.Info(req.OriginalSender)
	logger.Info("----------------------------")

	err, clientSender := SocketIdCacheFetch(req.OriginalSender)
	if err != nil {
		request.LogRequest(logger, req).Error("[worker_response_handler.go] Could not fetch original client Id: %s", req.OriginalSender)
		handler.errors = append(handler.errors, errors.New("Could not fetch original client Id"))
		return
	}

	handler.clientSender = clientSender
	return
}

func (handler *WorkerResponseHandler) relayResponse(req request.Request){
	if len(handler.errors) > 0 {
		request.LogRequest(logger, req).Error("[worker_response_handler.go] Unable to find orignal client sender Id: %s", handler.errors[0])
		return
	}

	msg, err := request.CreateMessage(req.ID,
		constants.WORKER_RESPONSE_RELAY,
		req.Data, req.Response, req.ServiceName, req.OriginalSender)

	if err != nil {
		request.LogRequest(logger, req).Error("[worker_response_handler.go] Unable to create message.")
		return
	}

	_, err = handler.broker.Socket.SendMessage(handler.clientSender, msg)

	if err != nil {
		request.LogRequest(logger, req).Error("[worker_response_handler.go] Unable to send WORKER_RESPONSE_RELAY message to client. ")
		logger.Error(fmt.Sprintf("[worker_response_handler.go] %s", err.Error()))
	}

	return
}
