package service

import (
	"github.com/supersid/iris2/request"
	"os"
	"runtime"
	"github.com/supersid/iris2/constants"
	simplLogger "github.com/supersid/iris2/logger"
	"github.com/Sirupsen/logrus"
	"fmt"
	"errors"
	"time"
)

var logger *logrus.Logger

type Service struct {
	ServiceName    string
	ClientRequest  []request.Request
	WaitingWorkers []*ServiceWorker
	Workers        map[string]*ServiceWorker
}

type ServiceWorker struct {
	Identity string
	Sender   string
	Expiry   time.Time
}

func NewServiceWorker(id, sender string) *ServiceWorker{
	return &ServiceWorker{
		Identity: id,
		Sender:   sender,
		Expiry:   time.Now().Add(constants.WORKER_HEARTBEAT_INTERVAL),
	}
}

func (s *Service) AddServiceWorker(sw *ServiceWorker){
	alreadyPresent := false
	for serviceId, _ := range s.Workers {
		if serviceId == sw.Identity {
			alreadyPresent = true
			break
		}
	}

	if !alreadyPresent{
		s.WaitingWorkers = append(s.WaitingWorkers, sw)
		s.Workers[sw.Identity] = sw
		logger.Info(fmt.Sprintf("[service.go] Added new service worker ID: %s to service %s", sw.Identity, s.ServiceName))
	}

}

func (s *Service) FindServiceWorker(serviceWorkerIdentity string) (*ServiceWorker){
	return s.Workers[serviceWorkerIdentity]
}

func (s *Service) AddClientRequest(req request.Request){
	logger.Info(fmt.Sprintf("[service.go] New Request from client ID %s for %s service", req.ID, req.ServiceName))
	s.ClientRequest = append(s.ClientRequest, req)
}

func (s *Service) ProcessRequests()(error, request.Request, *ServiceWorker){
	var req request.Request
	var sw *ServiceWorker
	var err error
	if len(s.ClientRequest) == 0 {
		logger.Info("No client requests to process.")
		err = errors.New("[service.go] No client requests to process.")

	}

	logger.Info(fmt.Sprintf("%d Waiting workers for %s service", len(s.WaitingWorkers), s.ServiceName))
	if len(s.WaitingWorkers) == 0 {
		logger.Info("No workers available to process requests.")
		err = errors.New("[service.go] No workers available to process requests.")
	}

	if err != nil {

		return err, request.Request{}, &ServiceWorker{}
	}


	req = s.PopFirstRequest()
	sw = s.PopFirstWorker()
	logger.Debug("[service.go] Waiting worker is %s", sw)
	return nil, req, sw
}

func (s *Service) PopFirstRequest()(request.Request){
	req := s.ClientRequest[0]
	s.ClientRequest = s.ClientRequest[1:]
	return req
}

func (s *Service) PopFirstWorker()(*ServiceWorker){
	worker := s.WaitingWorkers[0]
	s.WaitingWorkers = s.WaitingWorkers[1:]
	delete(s.Workers, worker.Identity)
	return worker
}

func (s *Service) DeleteWorker(sw *ServiceWorker){
	delete(s.Workers, sw.Identity)


	for i, sworker := range s.WaitingWorkers {
		if sworker.Identity == sw.Identity {
			s.WaitingWorkers = append(s.WaitingWorkers[:i], s.WaitingWorkers[i+1:]...)
			break
		}
	}

	logger.Info(fmt.Sprintf("[service.go] Deleting worker with ID: %s", sw.Identity))
}

func (s *ServiceWorker) UpdateExpiry(){
	s.Expiry.Add(constants.WORKER_HEARTBEAT_INTERVAL)
}

func (s *Service) HasNoWorkers() bool{
	if len(s.WaitingWorkers) == 0 {
		return true
	} else {
		return false
	}
}

func (s *Service) ClearAll(){
	s.ClientRequest = make([]request.Request, 0)
	s.WaitingWorkers = make([]*ServiceWorker, 0)
	s.Workers = make(map[string]*ServiceWorker)
	return
}

func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	env := os.Getenv("IRIS_ENV")

	if env == "" {
		env = constants.DEVELOPMENT_ENV
	}

	logger = simplLogger.Init(env, "")

}
