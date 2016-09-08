package worker

import (
	"github.com/supersid/iris2/constants"
)

func (worker *Worker) ProcessMessage(msg []string, messageChannel chan WorkerRequest){
	wr := NewWorkerRequest(msg)
	logger.Debug(wr)
	if wr.Command == constants.CLIENT_REQUEST_TO_WORKER {
		logger.Info("[worker.go] New client request to worker.")
		worker.UpdateHeartBeatExpiry()
		messageChannel <- wr
	} else if wr.Command == constants.BROKER_HEARTBEAT {
		worker.UpdateHeartBeatExpiry()
	}
}
