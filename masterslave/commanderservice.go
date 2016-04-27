package masterslave

import (
	"github.com/williammuji/shiran/shiran"
	"github.com/golang/glog"
)

type CommanderService struct {
	commander		*Commander
}

func NewCommanderService(commander *Commander) *CommanderService {
	service := &CommanderService{
		commander:		commander,
	}
	return service
}

func (service *CommanderService) HandleAddApplicationResponse(msg *AddApplicationResponse, session *shiran.Session) {
	glog.Infof("CommanderService HandleAddApplicationResponse %v", msg)
	service.commander.close()
}

func (service *CommanderService) HandleRemoveApplicationsResponse(msg *RemoveApplicationsResponse, session *shiran.Session) {
	glog.Infof("CommanderService HandleRemoveApplicationsResponse %v", msg)
	service.commander.close()
}

func (service *CommanderService) HandleStartApplicationsResponse(msg *StartApplicationsResponse, session *shiran.Session) {
	glog.Infof("CommanderService HandleStartApplicationsResponse %v", msg)
	service.commander.close()
}

func (service *CommanderService) HandleStopApplicationResponse(msg *StopApplicationResponse, session *shiran.Session) {
	glog.Infof("CommanderService HandleStopApplicationResponse %v", msg)
	service.commander.close()
}

func (service *CommanderService) HandleRestartApplicationResponse(msg *RestartApplicationResponse, session *shiran.Session) {
	glog.Infof("CommanderService HandleRestartApplicationResponse %v", msg)
	service.commander.close()
}

func (service *CommanderService) HandleListApplicationsResponse(msg *ListApplicationsResponse, session *shiran.Session) {
	glog.Infof("CommanderService HandleListApplicationsResponse %v", msg)
	service.commander.close()
}

func (service *CommanderService) HandleGetApplicationsResponse(msg *GetApplicationsResponse, session *shiran.Session) {
	glog.Infof("CommanderService HandleGetApplicationsResponse %v", msg)
	service.commander.close()
}
