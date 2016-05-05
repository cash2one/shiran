package client

import (
	"crypto/aes"
	"github.com/golang/glog"
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/userproto"
)

type GateClient struct {
	gateRegister			chan *shiran.Session
	gateUnregister			chan *shiran.Session
	client					*shiran.Client
	gateService				*GateService
	account					*Account
	key						[]byte
}

func NewGateClient(gateAddress string, account *Account, key []byte) *GateClient {
	gateClient := &GateClient{
		gateRegister:		make(chan *shiran.Session),
		gateUnregister:		make(chan *shiran.Session),
		account:			account,
		key:				key,
	}
	gateClient.client = shiran.NewClient(gateAddress, 1, gateClient.gateRegister, gateClient.gateUnregister, nil)
	gateClient.gateService = NewGateService(gateClient)
	return gateClient
}

func (gateClient *GateClient) Run() {
	gateClient.client.RegisterService(gateClient.gateService)
	gateClient.client.ConnectServer()
	gateClient.timer()
}

func (gateClient *GateClient) timer() {
	for {
		select {
		case session := <-gateClient.gateRegister:
			gateClient.onGateConnection(session)
		case _ = <-gateClient.gateUnregister:
			break
		}
	}
}

func (gateClient *GateClient) onGateConnection(session *shiran.Session) {
	c, err := aes.NewCipher(gateClient.key)
	if err != nil {
		glog.Errorf("GateClient onGateConnection NewCipher name:%s (%d bytes) = %s", gateClient.account.Name, len(gateClient.key), err)
		return
	}
	encryptedKey := make([]byte, len(gateClient.key))
	c.Encrypt(encryptedKey, gateClient.key)

	request := &UserLoginGateRequest{
		Name:				&gateClient.account.Name,
		EncryptedKey:		encryptedKey,
	}
	session.SendMessage("GateService", "HandleUserLoginGateRequest", request)
}
