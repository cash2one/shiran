package masterslave

import (
	"time"
	"runtime"
	"github.com/williammuji/shiran/shiran"
	"github.com/golang/glog"
)

type Master struct {
	//master
	masterRegister			chan *shiran.Session
	masterUnregister		chan *shiran.Session
	masterServer			*shiran.Server
	masterService			*MasterService
	masterSessions			map[string]*shiran.Session
	//command
	commandRegister			chan *shiran.Session
	commandUnregister		chan *shiran.Session
	commandServer			*shiran.Server
	commandSessions			map[string]*shiran.Session
	//config
	masterConfig			*MasterConfig
}

func NewMaster(masterConfig *MasterConfig) *Master {
	m := &Master{
		masterRegister:			make(chan *shiran.Session),
		masterUnregister:		make(chan *shiran.Session),
		masterSessions:			make(map[string]*shiran.Session),
		commandRegister:		make(chan *shiran.Session),
		commandUnregister:		make(chan *shiran.Session),
		commandSessions:		make(map[string]*shiran.Session),
		masterConfig:			masterConfig,
	}
	m.masterService = NewMasterService(masterConfig)
	m.masterServer = shiran.NewServer(m.masterRegister, m.masterUnregister, nil)
	m.commandServer = shiran.NewServer(m.commandRegister, m.commandUnregister, nil)
	return m 
}

func (m *Master) Run() {
	m.masterServer.RegisterService(m.masterService)
	m.commandServer.RegisterService(m.masterService)

	go m.timer()

	go m.commandServer.ListenAndServe(m.masterConfig.CommandAddress)

	m.masterServer.ListenAndServe(m.masterConfig.MasterAddress)
}

func (m *Master) timer() {
	ticks := time.Tick(time.Second * 1)

	for {
		select {
		case session := <-m.masterRegister:
			m.masterSessions[session.Name] = session
		case session := <-m.masterUnregister:
			msg := new(RemoveSlaveRequest)
			m.masterService.HandleRemoveSlaveRequest(msg, session)
			delete(m.masterSessions, session.Name)
		case session := <-m.commandRegister:
			msg := new(AddCommanderRequest)
			m.masterService.HandleAddCommanderRequest(msg, session)
			m.commandSessions[session.Name] = session
		case session := <-m.commandUnregister:
			msg := new(RemoveCommanderRequest)
			m.masterService.HandleRemoveCommanderRequest(msg, session)
			delete(m.commandSessions, session.Name)
		case _ = <-ticks:
			glog.Infof("masterSessions:%d commandSessions:%d %d", len(m.masterSessions), len(m.commandSessions), runtime.NumGoroutine())
		}
	}
}
