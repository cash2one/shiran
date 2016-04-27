package masterslave

import (
	"time"
	"runtime"
	"github.com/williammuji/shiran/shiran"
	"github.com/golang/glog"
)

type CommandOptions struct {
	ConfigFile		string
	SlaveName		string
	Add				bool
	Remove			bool
	Start			bool
	Stop			bool
	Restart			bool
	Get				bool
	List			bool
	AppName			string
	AppNames		[]string
}

type Commander struct {
	//master
	masterRegister			chan *shiran.Session
	masterUnregister		chan *shiran.Session
	masterClient			*shiran.Client
	session					*shiran.Session
	commanderService		*CommanderService
	//config
	commanderConfig			*CommanderConfig
	opt						*CommandOptions
	//quit
	quit					chan bool
}

func NewCommander(opt *CommandOptions) *Commander {
	c := &Commander{
		masterRegister:			make(chan *shiran.Session),
		masterUnregister:		make(chan *shiran.Session),
		opt:					opt,
		commanderConfig:		NewCommanderConfig(opt.ConfigFile),
		quit:					make(chan bool),
	}
	c.commanderService = NewCommanderService(c)
	c.masterClient = shiran.NewClient(c.commanderConfig.MasterAddress, 1, c.masterRegister, c.masterUnregister, nil)
	return c 
}

func (c *Commander) Run() {
	c.masterClient.RegisterService(c.commanderService)

	c.masterClient.ConnectServer()

	glog.Infof("run masterAddress:%s", c.commanderConfig.MasterAddress)

	c.timer()
}

func (c *Commander) timer() {
	ticks := time.Tick(time.Second * 1)

	for {
		select {
		case s := <-c.masterRegister:
			c.session = s
			c.onConnection()
		case _ = <-c.masterUnregister:
			c.session = nil
			return	
		case _ = <-c.quit:
			c.session.Close()
		case _ = <-ticks:
			glog.Infof("routines:%d", runtime.NumGoroutine())
		}
	}
}

func (c *Commander) onConnection() {
	if c.opt.Add == true {
		c.addApp(c.opt.SlaveName, c.opt.AppName)
	} else if c.opt.Remove == true {
		c.removeApps(c.opt.SlaveName, c.opt.AppNames)
	} else if c.opt.Start == true {
		c.startApps(c.opt.SlaveName, c.opt.AppNames)
	} else if c.opt.Stop == true {
		c.stopApp(c.opt.SlaveName, c.opt.AppName)
	} else if c.opt.Restart == true {
		c.restartApp(c.opt.SlaveName, c.opt.AppName)
	} else if c.opt.Get == true {
		c.getApps(c.opt.SlaveName, c.opt.AppNames)
	} else if c.opt.List == true {
		c.listApps(c.opt.SlaveName)
	}
}

func (c *Commander) addApp(slaveName, appName string) {
	msg := new(AddApplicationRequest)
	msg.Name = &appName
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	c.session.SendMessage("MasterService", "HandleAddApplicationRequest", msg)
}

func (c *Commander) removeApps(slaveName string, appNames []string) {
	msg := new(RemoveApplicationsRequest)
	msg.Names = appNames
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	c.session.SendMessage("MasterService", "HandleRemoveApplicationsRequest", msg)
}

func (c *Commander) startApps(slaveName string, appNames []string) {
	msg := new(StartApplicationsRequest)
	msg.Names = appNames
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	c.session.SendMessage("MasterService", "HandleStartApplicationsRequest", msg)
}

func (c *Commander) stopApp(slaveName, appName string) {
	msg := new(StopApplicationRequest)
	msg.Name = &appName
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	c.session.SendMessage("MasterService", "HandleStopApplicationRequest", msg)
}

func (c *Commander) restartApp(slaveName, appName string) {
	msg := new(RestartApplicationRequest)
	msg.Name = &appName
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	c.session.SendMessage("MasterService", "HandleRestartApplicationRequest", msg)
}

func (c *Commander) getApps(slaveName string, appNames[] string) {
	msg := new(GetApplicationsRequest)
	msg.Names = appNames
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	c.session.SendMessage("MasterService", "HandleGetApplicationsRequest", msg)
}

func (c *Commander) listApps(slaveName string) {
	msg := new(GetApplicationsRequest)
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	c.session.SendMessage("MasterService", "HandleListApplicationsRequest", msg)
}

func (c *Commander) close() {
	c.quit <- true
}
