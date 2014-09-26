package scalerservice

import (
	"fmt"

	"go-bro/broker"
	"go-bro/config"
)

type ScalerService struct {
	apps map[string]string
}

func New(appsBag map[string]string) *ScalerService {
	return &ScalerService{apps: appsBag}
}

func (s *ScalerService) Create(serviceInstance string, req broker.ServiceRequest, planConfig config.PlanConfig) error {
	return nil
}

func (s *ScalerService) Destroy(serviceInstance string) error {
	return nil
}

func (s *ScalerService) Bind(serviceInstance string, bindID string, req broker.BindRequest, planConfig config.PlanConfig) (broker.BindResponse, error) {
	fmt.Println("Adding APP", req.AppID, "with bindID:", bindID)
	s.apps[bindID] = req.AppID

	cred := broker.BindCredentials{}
	resp := broker.BindResponse{Credentials: cred}
	return resp, nil
}

func (s *ScalerService) Unbind(servceInstance string, bindID string) error {
	fmt.Println("Unbind id", bindID)
	delete(s.apps, bindID)
	return nil
}
