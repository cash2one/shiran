package masterslave

import (
	"time"
	"runtime"
	"github.com/williammuji/shiran/shiran"
	"github.com/golang/glog"
)

type Slave struct {
	//config
	slaveConfig				*SlaveConfig
	//master
	masterRegister			chan *shiran.Session
	masterUnregister		chan *shiran.Session
	masterClient			*shiran.Client
	slaveService			*SlaveService
	masterSessions			map[string]*shiran.Session
	//heartbeat
	heartbeats				map[string]*Heartbeat
}

func NewSlave(slaveConfig *SlaveConfig) *Slave {
	s := &Slave{
		slaveConfig:			slaveConfig,
		masterRegister:			make(chan *shiran.Session),
		masterUnregister:		make(chan *shiran.Session),
		slaveService:			NewSlaveService(slaveConfig.Name),
		masterSessions:			make(map[string]*shiran.Session),
		heartbeats:				make(map[string]*Heartbeat),
	}
	s.masterClient = shiran.NewClient(slaveConfig.MasterAddress, 1, s.masterRegister, s.masterUnregister, nil)
	return s 
}

func (s *Slave) Run() {
	s.masterClient.RegisterService(s.slaveService)

	s.masterClient.ConnectServer()

	glog.Infof("run slave %s masterAddress:%s", s.slaveConfig.Name, s.slaveConfig.MasterAddress)

	s.timer()
}

func (s *Slave) timer() {
	ticks := time.Tick(time.Second * 1)

	for {
		select {
		case session := <-s.masterRegister:
			s.masterSessions[session.Name] = session
			s.onConnection(session)

			heartbeat := NewHeartbeat(s.slaveConfig, session, s.slaveService)
			s.heartbeats[session.Name] = heartbeat
			heartbeat.start()
		case session := <-s.masterUnregister:
			delete(s.masterSessions, session.Name)
			if heartbeat, ok := s.heartbeats[session.Name]; ok {
				heartbeat.stop()
				delete(s.heartbeats, session.Name)
			}
		case _ = <-ticks:
			glog.Infof("masterSessions:%d routines:%d", len(s.masterSessions), runtime.NumGoroutine())
		}
	}
}

func (s *Slave) onConnection(session *shiran.Session) {
	if session.Closed == false {
		request := new(AddSlaveRequest)
		request.SlaveName = &s.slaveConfig.Name 
		session.SendMessage("MasterService", "HandleAddSlaveRequest", request)
	}
}
