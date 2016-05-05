package gate

import (
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/userproto"
	. "github.com/williammuji/shiran/proto/gatelogin"
	"github.com/golang/glog"
)

type GateService struct {
	userManager		*UserManager
}

func NewGateService(userManager *UserManager) *GateService {
	gs := &GateService{
		userManager:	userManager,
	}
	go gs.userManager.eventLoop()
	return gs
}

func (gs *GateService) HandleUserRandomKeyRequest(request *UserRandomKeyRequest, session *shiran.Session) {
	glog.Infof("%+v", request)

	event := userEvent{
		eventType:		USER_RANDOM_KEY_REQUEST,
		msg:			request,
		session:		session,
	}
	gs.userManager.PostEvent(event)
}

func (gs *GateService) HandleUserLoginGateRequest(request *UserLoginGateRequest, session *shiran.Session) {
	glog.Infof("%+v", request)

	event := userEvent{
		eventType:		USER_LOGIN_GATE_REQUEST,
		msg:			request,
		session:		session,
	}
	gs.userManager.PostEvent(event)
}
