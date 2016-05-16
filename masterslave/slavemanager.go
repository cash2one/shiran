package masterslave 

import (
	"github.com/golang/protobuf/proto"
	"github.com/williammuji/shiran/shiran"
	"github.com/golang/glog"
)

type SlaveEventType int

const (
	ADDSLAVE SlaveEventType = 1 << iota
	REMOVESLAVE

	ADDCOMMANDER
	REMOVECOMMANDER

	ADDAPP
	REMOVEAPPS
	STARTAPPS
	STOPAPP
	RESTARTAPP
	GETAPPS
	LISTAPPS

	ADDAPPRESPONSE
	REMOVEAPPSRESPONSE
	STARTAPPSRESPONSE
	STOPAPPRESPONSE
	RESTARTAPPRESPONSE
	GETAPPSRESPONSE
	LISTAPPSRESPONSE

	GETHARDWARE
	GETHARDWARERESPONSE

	GETFILECONTENT
	GETFILECONTENTRESPONSE

	GETFILECHECKSUM
	GETFILECHECKSUMRESPONSE

	RUNCOMMAND
	RUNCOMMANDRESPONSE
	RUNSCRIPT
)

type SlaveEvent struct {
	eventType       SlaveEventType
	msg             proto.Message   
	session         *shiran.Session
}

type SlaveManager struct {
	eventChannel		chan SlaveEvent
	slaves				map[string]*shiran.Session
	commanders			map[string]*shiran.Session
	masterConfig		*MasterConfig
}

func NewSlaveManager(masterConfig *MasterConfig) *SlaveManager {
	am := &SlaveManager{
		eventChannel:		make(chan SlaveEvent, 10),	//FIXME 10
		slaves:				make(map[string]*shiran.Session),
		commanders:			make(map[string]*shiran.Session),
		masterConfig:		masterConfig,
	}
	return am
}

func (am *SlaveManager) PostEvent(event SlaveEvent) {
	am.eventChannel <- event
}

func (am *SlaveManager) eventLoop() {
	for event := range am.eventChannel {
		switch event.eventType {
			//commander
		case ADDCOMMANDER:
			if _, ok := event.msg.(*AddCommanderRequest); ok {
				am.commanders[event.session.Name] = event.session
			}
		case REMOVECOMMANDER:
			if _, ok := event.msg.(*RemoveCommanderRequest); ok {
				for name, s := range am.commanders {
					if s == event.session {
						delete(am.commanders, name)
						break
					}
				}
			}
			//slave
		case ADDSLAVE:
			if request, ok := event.msg.(*AddSlaveRequest); ok {
				am.slaves[request.GetSlaveName()] = event.session
			}
		case REMOVESLAVE:
			if _, ok := event.msg.(*RemoveSlaveRequest); ok {
				for name, s := range am.slaves {
					if s == event.session {
						delete(am.slaves, name)
						break
					}
				}
			}
			//app request
		case ADDAPP:
			if request, ok := event.msg.(*AddApplicationRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name

					for _, slave := range am.masterConfig.Slave {
						if slave.Name == request.GetSlaveCommander().GetSlaveName() {
							for _, app := range slave.App {
								if app.Name == request.GetName() {
									request.Binary = &app.Bin
									request.Args = app.Arg
									break
								}
							}
						}
					}

					session.SendMessage("SlaveService", "HandleAddApplicationRequest", request)
				} else {
					glog.Errorf("eventLoop ADDAPP slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case REMOVEAPPS:
			if request, ok := event.msg.(*RemoveApplicationsRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleRemoveApplicationsRequest", request)
				} else {
					glog.Errorf("eventLoop REMOVEAPPS slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case STARTAPPS:
			if request, ok := event.msg.(*StartApplicationsRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleStartApplicationsRequest", request)
				} else {
					glog.Errorf("eventLoop STARTAPP slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case STOPAPP:
			if request, ok := event.msg.(*StopApplicationRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleStopApplicationRequest", request)
				} else {
					glog.Errorf("eventLoop STOPAPP slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case RESTARTAPP:
			if request, ok := event.msg.(*RestartApplicationRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleRestartApplicationRequest", request)
				} else {
					glog.Errorf("eventLoop RESTARTAPP slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case GETAPPS:
			if request, ok := event.msg.(*GetApplicationsRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleGetApplicationsRequest", request)
				} else {
					glog.Errorf("eventLoop GETAPPS slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case LISTAPPS:
			if request, ok := event.msg.(*ListApplicationsRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleListApplicationsRequest", request)
				} else {
					glog.Errorf("eventLoop LISTAPPS slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case GETHARDWARE:
			if request, ok := event.msg.(*GetHardwareRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleGetHardwareRequest", request)
				} else {
					glog.Errorf("eventLoop GETHARDWARE slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case GETFILECONTENT:
			if request, ok := event.msg.(*GetFileContentRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleGetFileContentRequest", request)
				} else {
					glog.Errorf("eventLoop GETFILECONTENT slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case GETFILECHECKSUM:
			if request, ok := event.msg.(*GetFileChecksumRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleGetFileChecksumRequest", request)
				} else {
					glog.Errorf("eventLoop GETFILECHECKSUM slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case RUNCOMMAND:
			if request, ok := event.msg.(*RunCommandRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleRunCommandRequest", request)
				} else {
					glog.Errorf("eventLoop RUNCOMMAND slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
		case RUNSCRIPT:
			if request, ok := event.msg.(*RunScriptRequest); ok {
				if session, exist := am.slaves[request.GetSlaveCommander().GetSlaveName()]; exist {
					request.SlaveCommander.CommanderName = &event.session.Name
					session.SendMessage("SlaveService", "HandleRunScriptRequest", request)
				} else {
					glog.Errorf("eventLoop RUNSCRIPT slaveName:%s not exist", request.GetSlaveCommander().GetSlaveName())
				}
			}
			//app response
		case ADDAPPRESPONSE:
			if response, ok := event.msg.(*AddApplicationResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleAddApplicationResponse", response)
				} else {
					glog.Errorf("eventLoop ADDAPPRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case STARTAPPSRESPONSE:
			if response, ok := event.msg.(*StartApplicationsResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleStartApplicationsResponse", response)
				} else {
					glog.Errorf("eventLoop STARTAPPSRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case STOPAPPRESPONSE:
			if response, ok := event.msg.(*StopApplicationResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleStopApplicationResponse", response)
				} else {
					glog.Errorf("eventLoop STOPAPPRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case RESTARTAPPRESPONSE:
			if response, ok := event.msg.(*RestartApplicationResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleRestartApplicationResponse", response)
				} else {
					glog.Errorf("eventLoop RESTARTAPPRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case LISTAPPSRESPONSE:
			if response, ok := event.msg.(*ListApplicationsResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleListApplicationsResponse", response)
				} else {
					glog.Errorf("eventLoop LISTAPPSRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case GETAPPSRESPONSE:
			if response, ok := event.msg.(*GetApplicationsResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleGetApplicationsResponse", response)
				} else {
					glog.Errorf("eventLoop GETAPPSRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case REMOVEAPPSRESPONSE:
			if response, ok := event.msg.(*RemoveApplicationsResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleRemoveApplicationsResponse", response)
				} else {
					glog.Errorf("eventLoop REMOVEAPPSRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case GETHARDWARERESPONSE:
			if response, ok := event.msg.(*GetHardwareResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleGetHardwareResponse", response)
				} else {
					glog.Errorf("eventLoop GETHARDWARERESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case GETFILECONTENTRESPONSE:
			if response, ok := event.msg.(*GetFileContentResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleGetFileContentResponse", response)
				} else {
					glog.Errorf("eventLoop GETFILECONTENTRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case GETFILECHECKSUMRESPONSE:
			if response, ok := event.msg.(*GetFileChecksumResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleGetFileChecksumResponse", response)
				} else {
					glog.Errorf("eventLoop GETFILECHECKSUMRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		case RUNCOMMANDRESPONSE:
			if response, ok := event.msg.(*RunCommandResponse); ok {
				if session, exist := am.commanders[response.GetSlaveCommander().GetCommanderName()]; exist {
					session.SendMessage("CommanderService", "HandleRunCommandResponse", response)
				} else {
					glog.Errorf("eventLoop RUNCOMMANDRESPONSE commanderName:%s not exist", response.GetSlaveCommander().GetCommanderName())
				}
			}
		}
	}
}
