package client

import (
	"github.com/golang/glog"
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/userproto"
)

type LoginClient struct {
	loginRegister			chan *shiran.Session
	loginUnregister			chan *shiran.Session
	client					*shiran.Client
	loginService			*LoginService
	account					*Account
}

func NewLoginClient(loginAddress string, account *Account) *LoginClient {
	loginClient := &LoginClient{
		loginRegister:		make(chan *shiran.Session),
		loginUnregister:	make(chan *shiran.Session),
		account:			account,
	}
	loginClient.client = shiran.NewClient(loginAddress, 1, loginClient.loginRegister, loginClient.loginUnregister, nil)
	loginClient.loginService = NewLoginService(loginClient)
	return loginClient
}

func (loginClient *LoginClient) Run(caFile string) {
	loginClient.client.RegisterService(loginClient.loginService)
	tlsConfig := shiran.GetClientTlsConfiguration(caFile)
	loginClient.client.TlsConnectServer(tlsConfig)
	loginClient.timer()
}

func (loginClient *LoginClient) timer() {
	for {
		select {
		case session := <-loginClient.loginRegister:
			loginClient.onLoginConnection(session)
		case _ = <-loginClient.loginUnregister:
			break
		}
	}
}

func (loginClient *LoginClient) onLoginConnection(session *shiran.Session) {
	request := &UserLoginLoginRequest{
		Name:		&loginClient.account.Name,
		Passwd:		&loginClient.account.Passwd,
		Zone:		&loginClient.account.Zone,
	}
	glog.Errorf("%v", request)
	session.SendMessage("LoginService", "HandleUserLoginLoginRequest", request)
}
