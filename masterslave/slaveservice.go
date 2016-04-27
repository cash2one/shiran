package masterslave

import (
	"github.com/williammuji/shiran/shiran"
)

type SlaveService struct {
	appManager	*AppManager
}

func NewSlaveService() *SlaveService {
	service := &SlaveService{
		appManager:		NewAppManager(),
	}
	go service.appManager.eventLoop()
	return service
}

func (service *SlaveService) HandleAddApplicationRequest(msg *AddApplicationRequest, session *shiran.Session) {
	event := Event{
		eventType:		ADD,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleRemoveApplicationsRequest(msg *RemoveApplicationsRequest, session *shiran.Session) {
	event := Event{
		eventType:		REMOVE,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleStartApplicationsRequest(msg *StartApplicationsRequest, session *shiran.Session) {
	event := Event{
		eventType:		START,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleStopApplicationRequest(msg *StopApplicationRequest, session *shiran.Session) {
	event := Event{
		eventType:		STOP,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleRestartApplicationRequest(msg *RestartApplicationRequest, session *shiran.Session) {
	event := Event{
		eventType:		RESTART,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleListApplicationsRequest(msg *ListApplicationsRequest, session *shiran.Session) {
	event := Event{
		eventType:		LIST,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}

func (service *SlaveService) HandleGetApplicationsRequest(msg *GetApplicationsRequest, session *shiran.Session) {
	event := Event{
		eventType:		GET,
		msg:			msg,
		session:		session,
	}
	service.appManager.PostEvent(event)
}
