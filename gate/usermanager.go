package gate

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/userproto"
	. "github.com/williammuji/shiran/proto/gatelogin"
	"crypto/aes"
	"bytes"
)

type userEventType int

const (
	USER_LOGIN_GATE_REQUEST = 1 << iota
	USER_UNREGISTER
	USER_RANDOM_KEY_REQUEST
)

type userEvent struct {
	eventType	userEventType
	msg			proto.Message
	session		*shiran.Session
}

type UserManager struct {
	eventChannel		chan userEvent
	sessions			map[string]*shiran.Session
	randomKeys			map[string][]byte
}

func NewUserManager() *UserManager {
	um := &UserManager{
		eventChannel:			make(chan userEvent, 1024),
		sessions:				make(map[string]*shiran.Session),
		randomKeys:				make(map[string][]byte),
	}
	return um
}

func (um *UserManager) eventLoop() {
	for event := range um.eventChannel {
		switch event.eventType {
		case USER_RANDOM_KEY_REQUEST:
			if request, ok := event.msg.(*UserRandomKeyRequest); ok {
				um.randomKeys[request.GetName()] = request.GetRandomKey()		//is it null
			}
		case USER_LOGIN_GATE_REQUEST:
			if request, ok := event.msg.(*UserLoginGateRequest); ok {
				response := &UserLoginGateResponse{}
				if len(request.GetEncryptedKey()) != 16 {
					response.State = UserLoginGateState_kEncryptError.Enum()
					event.session.SendMessage("GateService", "HandleUserLoginGateResponse", response)
					glog.Errorf("UserManager eventLoop recv EncryptedKey len != 16 %s", request.GetName())
					return
				}

				if key, exist := um.randomKeys[request.GetName()]; !exist {
					response.State = UserLoginGateState_kNameError.Enum()
					event.session.SendMessage("GateService", "HandleUserLoginGateResponse", response)
				} else {
					c, err := aes.NewCipher(key)
					if err != nil {
						response.State = UserLoginGateState_kEncryptError.Enum()
						event.session.SendMessage("GateService", "HandleUserLoginGateResponse", response)
						glog.Errorf("UserManager eventLoop NewCipher name:%s (%d bytes) = %s", request.GetName(), len(key), err)
						return
					}
					encryptedKey := make([]byte, len(key))
					c.Encrypt(encryptedKey, key)
					recvEncryptedKey := request.GetEncryptedKey()
					if !bytes.Equal(encryptedKey, recvEncryptedKey) {
						response.State = UserLoginGateState_kEncryptError.Enum()
						event.session.SendMessage("GateService", "HandleUserLoginGateResponse", response)
						glog.Errorf("UserManager eventLoop encryptedKey not equal name:%s %v", request.GetName(), encryptedKey)
						return
					}

					um.sessions[request.GetName()] = event.session
					response.State = UserLoginGateState_kLoginGateSuccess.Enum()
					event.session.SendMessage("GateService", "HandleUserLoginGateResponse", response)
				}
			}
		case USER_UNREGISTER:
			var name string
			var s *shiran.Session
			for name, s = range um.sessions {
				if s == event.session {
					delete(um.sessions, name)
					break
				}
			}
			delete(um.randomKeys, name)
		}
	}
}

func (um *UserManager) PostEvent(event userEvent) {
	um.eventChannel <- event
}
