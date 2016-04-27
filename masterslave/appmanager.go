package masterslave 

import (
	"os"
	"os/exec"
	"github.com/golang/protobuf/proto"
	"github.com/williammuji/shiran/shiran"
	"github.com/golang/glog"
)

type Application struct {
	request			*AddApplicationRequest	
	status			*ApplicationStatus
	cmd				*exec.Cmd
}

type EventType int

const (
	//external
	ADD EventType = 1 << iota
	REMOVE
	START
	STOP
	RESTART
	LIST
	GET
	//internal
	DEAD
)

type Event struct {
	eventType		EventType
	msg				proto.Message	
	session			*shiran.Session
}

type AppManager struct {
	eventChannel		chan Event
	apps				map[string]*Application
}

func NewAppManager() *AppManager {
	am := &AppManager{
		eventChannel:		make(chan Event, 10),	//FIXME 10
		apps:				make(map[string]*Application),
	}
	return am
}

func (am *AppManager) PostEvent(event Event) {
	am.eventChannel <- event
}

func (am *AppManager) eventLoop() {
	for event := range am.eventChannel {
		switch event.eventType {
		case ADD:
			if request, ok := event.msg.(*AddApplicationRequest); ok {
				am.addApp(request, event.session)
			}
		case REMOVE:
			if request, ok := event.msg.(*RemoveApplicationsRequest); ok {
				am.removeApps(request, event.session)
			}
		case START:
			if request, ok := event.msg.(*StartApplicationsRequest); ok {
				am.startApps(request, event.session)
			}
		case STOP:
			if request, ok := event.msg.(*StopApplicationRequest); ok {
				am.stopApp(request, event.session)
			}
		case RESTART:
			if request, ok := event.msg.(*RestartApplicationRequest); ok {
				am.restartApp(request, event.session)
			}
		case LIST:
			if request, ok := event.msg.(*ListApplicationsRequest); ok {
				am.listApps(request, event.session)
			}
		case GET:
			if request, ok := event.msg.(*GetApplicationsRequest); ok {
				am.getApps(request, event.session)
			}
		case DEAD:
			if request, ok := event.msg.(*DeadApplicationRequest); ok {
				am.deadApp(request, event.session)
			}
		}
	}
}

func (am *AppManager) addApp(request *AddApplicationRequest, session *shiran.Session) {
	if app, ok := am.apps[request.GetName()]; ok {
		prevRequest := app.request
		app.request = request
		status := app.status
		status.Name = request.Name

		response := new(AddApplicationResponse)
		response.Status = status
		response.PrevRequest = prevRequest
		response.SlaveCommander = request.GetSlaveCommander()
		session.SendMessage("MasterService", "HandleAddApplicationResponse", response)
		glog.Info("addApp app:%s state:%s", request.GetName(), app.status.GetState())
	} else {
		status := new(ApplicationStatus)
		status.State = ApplicationState_kNewApp.Enum()
		status.Name = request.Name

		app := &Application{
			request:	request,
			status:		status,
		}
		am.apps[request.GetName()] = app
		response := &AddApplicationResponse{
			Status:			status,
			SlaveCommander:	request.GetSlaveCommander(),
		}
		session.SendMessage("MasterService", "HandleAddApplicationResponse", response)
		glog.Infof("addApp app:%s state:%s", request.GetName(), app.status.GetState())
	}
}

func (am *AppManager) removeApps(request *RemoveApplicationsRequest, session *shiran.Session) {
	response := &RemoveApplicationsResponse{
			SlaveCommander:	request.GetSlaveCommander(),
	}
	for _, name := range request.GetNames() {
		if app, ok := am.apps[name]; ok {
			glog.Infof("removeApp app:%s state:%s", app.status.GetName(), app.status.GetState())
			delete(am.apps, name)
		}
	}
	session.SendMessage("MasterService", "HandleRemoveApplicationsResponse", response)
}

func (am *AppManager) startApps(request *StartApplicationsRequest, session *shiran.Session) {
	response := &StartApplicationsResponse{
			SlaveCommander:	request.GetSlaveCommander(),
	}
	for index, name := range request.GetNames() {
		if app, ok := am.apps[name]; ok {
			if app.status.GetState() != ApplicationState_kRunning && app.status.GetState() != ApplicationState_kStopping && app.status.GetState() != ApplicationState_kRestarting {
				app.start(am, session)
				response.Status = append(response.Status, app.status)
			} else {
				response.Status = append(response.Status, app.status)
				response.Status[index].State = ApplicationState_kError.Enum()
			}
		} else {
			status := &ApplicationStatus{}
			status.State = ApplicationState_kUnknown.Enum()
			status.Name = &name
			m := string("Application is unknown.")
			status.Message = &m
			response.Status = append(response.Status, status)
		}
	}
	session.SendMessage("MasterService", "HandleStartApplicationsResponse", response)
}

func (am *AppManager) stopApp(request *StopApplicationRequest, session *shiran.Session) {
	if app, ok := am.apps[request.GetName()]; ok {
		if app.status.GetState() == ApplicationState_kRunning {
			app.status.State = ApplicationState_kStopping.Enum()
			app.cmd.Process.Kill()
			//syscall.Kill(int(app.status.GetPid()), syscall.SIGINT)
			response := &StopApplicationResponse{
				Status:			app.status,
				SlaveCommander:	request.GetSlaveCommander(),
			}
			session.SendMessage("MasterService", "HandleStopApplicationResponse", response)
		} else {
			glog.Errorf("stopApp app:%s not running", request.GetName())
		}
	} else {
		glog.Errorf("stopApp app:%s not found", request.GetName())
	}
}

func (am *AppManager) restartApp(request *RestartApplicationRequest, session *shiran.Session) {
	if app, ok := am.apps[request.GetName()]; ok {
		if app.status.GetState() == ApplicationState_kRunning {
			app.status.State = ApplicationState_kRestarting.Enum()
			app.cmd.Process.Kill()
			//syscall.Kill(int(app.status.GetPid()), syscall.SIGINT)
			response := &RestartApplicationResponse{
				Status:			app.status,
				SlaveCommander:	request.GetSlaveCommander(),
			}
			session.SendMessage("MasterService", "HandleRestartApplicationResponse", response)
		} else {
			glog.Errorf("restartApp app:%s not running", request.GetName())
		}
	} else {
		glog.Errorf("restartApp app:%s not found", request.GetName())
	}
}

func (am *AppManager) getApps(request *GetApplicationsRequest, session *shiran.Session) {
	response := &GetApplicationsResponse{
		SlaveCommander:	request.GetSlaveCommander(),
	}
	for _, name := range request.GetNames() {
		if app, ok := am.apps[name]; ok {
			response.Status = append(response.Status, app.status)
		}
	}
	session.SendMessage("MasterService", "HandleGetApplicationsResponse", response)
}

func (am *AppManager) listApps(request *ListApplicationsRequest, session *shiran.Session) {
	response := &ListApplicationsResponse{
		SlaveCommander:	request.GetSlaveCommander(),
	}
	for _, app := range am.apps {
		response.Names = append(response.Names, app.status.GetName())
	}
	session.SendMessage("MasterService", "HandleListApplicationsResponse", response)
}

func (am *AppManager) deadApp(request *DeadApplicationRequest, session *shiran.Session) {
	if app, ok := am.apps[request.GetAppName()]; ok {
		if app.status.GetPid() == request.GetProcessID() {
			app.status.State = ApplicationState_kExited.Enum()
			glog.Errorf("deadApp app:%s pid:%d processState:%s", request.GetAppName(), request.GetProcessID(), request.GetProcessState())

			if request.GetPrevState() == ApplicationState_kStopping {
				response := &StopApplicationResponse{
					SlaveCommander:	request.GetSlaveCommander(),
					Status:			app.status,
				}
				session.SendMessage("MasterService", "HandleStopApplicationResponse", response)
			} else if request.GetPrevState() == ApplicationState_kRestarting {
				response := &RestartApplicationResponse{
					SlaveCommander:	request.GetSlaveCommander(),
					Status:			app.status,
				}
				session.SendMessage("MasterService", "HandleRestartApplicationResponse", response)
				app.start(am, session)
				response.Status = app.status
				session.SendMessage("MasterService", "HandleRestartApplicationResponse", response)
			}
		} else {
			app.status.State = ApplicationState_kError.Enum()
			glog.Errorf("deadApp app:%s pid:%d processState:%s pid not found", request.GetAppName(), request.GetProcessID(), request.GetProcessState())
		}
	} else {
		glog.Errorf("deadApp app:%s pid:%d processState:%s not exist in apps", request.GetAppName(), request.GetProcessID(), request.GetProcessState())
	}
}

func (app *Application) start(am *AppManager, session *shiran.Session) {
	prevState := app.status.GetState()
	app.cmd = exec.Command(app.request.GetBinary(), app.request.GetArgs() ...)
	err := app.cmd.Start()
	if err != nil {
		app.status.State = ApplicationState_kError.Enum()
		glog.Errorf("start app:%s failed error:%s", app.status.GetName(), err)
		return
	} else {
		app.status.State = ApplicationState_kRunning.Enum()
		id := int32(app.cmd.Process.Pid)
		app.status.Pid = &id
		glog.Infof("start app:%s success pid:%d Process.Pid:%d Getpid:%d Getppid:%d state:%s args:%v", app.status.GetName(), app.status.GetPid(), app.cmd.Process.Pid, os.Getpid(), os.Getppid(), app.status.GetState(), app.cmd.Args)
	}

	go func(prevState ApplicationState, app *Application) {
		glog.Infof("start app:%s Waiting for command to finish...", app.status.GetName())
		err := app.cmd.Wait()
		glog.Infof("start app:%s Command finished with error: %v", app.status.GetName(), err)

		deadRequest := &DeadApplicationRequest{}
		deadRequest.AppName = app.status.Name
		deadRequest.PrevState = &prevState
		id := int32(app.cmd.ProcessState.Pid())
		deadRequest.ProcessID = &id
		state := app.cmd.ProcessState.String()
		deadRequest.ProcessState = &state
		am.PostEvent(Event{
			eventType:		DEAD,
			msg:			deadRequest,
			session:		session,
		})
	}(prevState, app)
}
