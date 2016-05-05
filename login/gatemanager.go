package login

import (
	"github.com/golang/protobuf/proto"
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/gatelogin"
	. "github.com/williammuji/shiran/proto/userproto"
)

type gateEventType int

const (
	ZONE_GATE_SESSION_REGISTER = 1 << iota
	ZONE_GATE_SESSION_UNREGISTER
	USER_LOGIN_REQUEST
)

type gateEvent struct {
	eventType	gateEventType
	msg			proto.Message
	session		*shiran.Session
}

type zoneSession struct {
	session				*shiran.Session
	gateListenAddress	string
}

type GateManager struct {
	eventChannel		chan gateEvent
	zoneSessions		map[int32][]*zoneSession
	nextID				int
}

func NewGateManager() *GateManager {
	gm := &GateManager{
		eventChannel:			make(chan gateEvent, 1024),
		zoneSessions:			make(map[int32][]*zoneSession),
	}
	return gm
}

func (gm *GateManager) eventLoop() {
	for event := range gm.eventChannel {
		switch event.eventType {
		case ZONE_GATE_SESSION_REGISTER:
			if request, ok := event.msg.(*GateLoginRequest); ok {
				zs := &zoneSession{
					session:				event.session,
					gateListenAddress:		request.GetGateListenAddress(),
				}
				gm.zoneSessions[request.GetZone()] = append(gm.zoneSessions[request.GetZone()], zs)
			}
		case ZONE_GATE_SESSION_UNREGISTER:
			for _, zs := range gm.zoneSessions {
				for index, s := range zs {
					if s.session == event.session {
						zs = append(zs[:index], zs[index+1:] ...)
						break
					}
				}
			}
		case USER_LOGIN_REQUEST:
			if request, ok := event.msg.(*UserRandomKeyRequest); ok {
				if zs, exist := gm.zoneSessions[request.GetZone()]; exist {
					if len(zs) > 0 {
						response := &UserLoginLoginResponse{}
						response.State = UserLoginLoginState_kLoginLoginSuccess.Enum()
						response.RandomKey = request.RandomKey
						response.GateListenAddress = &zs[gm.nextID].gateListenAddress
						event.session.SendMessage("LoginService", "HandleUserLoginLoginResponse", response)

						zs[gm.nextID].session.SendMessage("GateService", "HandleUserRandomKeyRequest", event.msg)
						gm.nextID++
						if gm.nextID >= len(zs) {
							gm.nextID = 0
						}
					}
				}
			}
		}
	}
}

func (gm *GateManager) PostEvent(event gateEvent) {
	gm.eventChannel <- event
}
