package masterslave

import (
	"github.com/williammuji/shiran/shiran"
	"github.com/golang/glog"
)

type MasterService struct {
	slaveManager			*SlaveManager
}

func NewMasterService(masterConfig *MasterConfig) *MasterService {
	service := &MasterService{
		slaveManager:		NewSlaveManager(masterConfig),
	}
	go service.slaveManager.eventLoop()
	return service
}

func (service *MasterService) HandleAddSlaveRequest(msg *AddSlaveRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		ADDSLAVE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRemoveSlaveRequest(msg *RemoveSlaveRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		REMOVESLAVE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}


func (service *MasterService) HandleAddCommanderRequest(msg *AddCommanderRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		ADDCOMMANDER,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRemoveCommanderRequest(msg *RemoveCommanderRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		REMOVECOMMANDER,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleAddApplicationRequest(msg *AddApplicationRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		ADDAPP,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRemoveApplicationsRequest(msg *RemoveApplicationsRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		REMOVEAPPS,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleStartApplicationsRequest(msg *StartApplicationsRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		STARTAPPS,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleStopApplicationRequest(msg *StopApplicationRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		STOPAPP,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRestartApplicationRequest(msg *RestartApplicationRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		RESTARTAPP,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleGetApplicationsRequest(msg *GetApplicationsRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		GETAPPS,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleListApplicationsRequest(msg *ListApplicationsRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		LISTAPPS,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleGetHardwareRequest(msg *GetHardwareRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		GETHARDWARE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleGetFileContentRequest(msg *GetFileContentRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		GETFILECONTENT,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleGetFileChecksumRequest(msg *GetFileChecksumRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		GETFILECHECKSUM,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRunCommandRequest(msg *RunCommandRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		RUNCOMMAND,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRunScriptRequest(msg *RunScriptRequest, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		RUNSCRIPT,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}



func (service *MasterService) HandleAddApplicationResponse(msg *AddApplicationResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		ADDAPPRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleStartApplicationsResponse(msg *StartApplicationsResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		STARTAPPSRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleStopApplicationResponse(msg *StopApplicationResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		STOPAPPRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRestartApplicationResponse(msg *RestartApplicationResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		RESTARTAPPRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleListApplicationsResponse(msg *ListApplicationsResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		LISTAPPSRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleGetApplicationsResponse(msg *GetApplicationsResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		GETAPPSRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRemoveApplicationsResponse(msg *RemoveApplicationsResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		REMOVEAPPSRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleGetHardwareResponse(msg *GetHardwareResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		GETHARDWARERESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleGetFileContentResponse(msg *GetFileContentResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		GETFILECONTENTRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleGetFileChecksumResponse(msg *GetFileChecksumResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		GETFILECHECKSUMRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}

func (service *MasterService) HandleRunCommandResponse(msg *RunCommandResponse, session *shiran.Session) {
	event := SlaveEvent{
		eventType:		RUNCOMMANDRESPONSE,
		msg:			msg,
		session:		session,
	}
	service.slaveManager.PostEvent(event)
}


func (service *MasterService) HandleSlaveHeartbeat(msg *SlaveHeartbeat, session *shiran.Session) {
	glog.Infof("session.Name:%s Heartbeat: %+v", session.Name, msg)
}
