package login

import (
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/gatelogin"
	"github.com/golang/glog"
)

type GateService struct {
	gateManager		*GateManager
}

func NewGateService(gateManager *GateManager) *GateService {
	gs := &GateService{
		gateManager:	gateManager,
	}
	return gs
}

func (gs *GateService) HandleGateLoginRequest(request *GateLoginRequest, session *shiran.Session) {
	glog.Infof("%v", request)

	event := gateEvent{
		eventType:		ZONE_GATE_SESSION_REGISTER,
		msg:			request,
		session:		session,
	}
	gs.gateManager.PostEvent(event)
}
