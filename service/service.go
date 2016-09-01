package service

import (
	"github.com/supersid/iris2/request"
)

type Service struct {
	ServiceName    string
	ClientRequest  []request.Request
	WaitingWorkers []*ServiceWorker
	Workers        map[string]*ServiceWorker
}

type ServiceWorker struct {
	Identity string
	Sender   string
}

func NewServiceWorker(id, sender string) *ServiceWorker{
	return &ServiceWorker{
		Identity: id,
		Sender:   sender,
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
	}

}
