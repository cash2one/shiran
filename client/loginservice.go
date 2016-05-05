package client 

import (
	"github.com/golang/glog"
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/userproto"
)

type LoginService struct {
	loginClient		*LoginClient
}

func NewLoginService(loginClient *LoginClient) *LoginService {
	service := &LoginService{
		loginClient:		loginClient,
	}
	return service
}

func (service *LoginService) HandleUserLoginLoginResponse(msg *UserLoginLoginResponse, session *shiran.Session) {
	glog.Errorf("%s %s %+v", session.Name, service.loginClient.account.Name, msg)

	if msg.GetState() == UserLoginLoginState_kLoginLoginSuccess {
		gateClient := NewGateClient(msg.GetGateListenAddress(), service.loginClient.account, msg.GetRandomKey())
		go gateClient.Run()
	}

	session.Close()
}
