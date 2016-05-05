package client 

import (
	"github.com/williammuji/shiran/shiran"
	"github.com/golang/glog"
	. "github.com/williammuji/shiran/proto/userproto"
)

type GateService struct {
	gateClient		*GateClient
}

func NewGateService(gateClient *GateClient) *GateService {
	service := &GateService{
		gateClient:		gateClient,
	}
	return service
}

func (service *GateService) HandleUserLoginGateResponse(msg *UserLoginGateResponse, session *shiran.Session) {
	glog.Infof("%+v", msg)

	session.Close()
}
