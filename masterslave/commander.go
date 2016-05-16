package masterslave

import (
	"time"
	"io/ioutil"
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
	GetHardware		bool
	Lshw			bool
	GetFileContent	bool
	FileName		string
	MaxSize			int64
	GetFileChecksum	bool
	Files			[]string
	RunCommand		bool
	Command			string
	Args			[]string
	MaxStdout		int
	MaxStderr		int
	Timeout			int
	RunScript		bool
	Script			string	
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
	} else if c.opt.GetHardware == true {
		c.getHardware(c.opt.SlaveName, c.opt.Lshw)
	} else if c.opt.GetFileContent == true {
		c.getFileContent(c.opt.SlaveName, c.opt.FileName, c.opt.MaxSize)
	} else if c.opt.GetFileChecksum == true {
		c.getFileChecksum(c.opt.SlaveName, c.opt.Files)
	} else if c.opt.RunCommand == true {
		c.runCommand(c.opt.SlaveName, c.opt.Command, c.opt.Args, int32(c.opt.MaxStdout), int32(c.opt.MaxStderr), int32(c.opt.Timeout))
	} else if c.opt.RunScript == true {
		c.runScript(c.opt.SlaveName, c.opt.Script, int32(c.opt.MaxStdout), int32(c.opt.MaxStderr), int32(c.opt.Timeout))
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

func (c *Commander) getApps(slaveName string, appNames []string) {
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

func (c *Commander) getHardware(slaveName string, lshw bool) {
	msg := new(GetHardwareRequest)
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	msg.Lshw = &lshw
	c.session.SendMessage("MasterService", "HandleGetHardwareRequest", msg)
}

func (c *Commander) getFileContent(slaveName, fileName string, maxSize int64) {
	msg := new(GetFileContentRequest)
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	msg.FileName = &fileName
	msg.MaxSize = &maxSize
	c.session.SendMessage("MasterService", "HandleGetFileContentRequest", msg)
}

func (c *Commander) getFileChecksum(slaveName string, fileNames []string) {
	msg := new(GetFileChecksumRequest)
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	msg.Files = fileNames
	c.session.SendMessage("MasterService", "HandleGetFileChecksumRequest", msg)
}

func (c *Commander) runCommand(slaveName string, command string, args []string, maxstdout, maxstderr int32, timeout int32) {
	msg := new(RunCommandRequest)
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName
	msg.Command = &command
	msg.Args = args
	msg.MaxStdout = &maxstdout
	msg.MaxStderr = &maxstderr
	msg.Timeout = &timeout
	c.session.SendMessage("MasterService", "HandleRunCommandRequest", msg)
}

func (c *Commander) runScript(slaveName string, script string, maxstdout, maxstderr int32, timeout int32) {
	msg := new(RunScriptRequest)
	msg.SlaveCommander = new(SlaveCommander)
	msg.SlaveCommander.SlaveName = &slaveName

	buffer, fileErr := ioutil.ReadFile(script)
	if fileErr != nil {
		glog.Errorf("load script %s failed %s", script, fileErr)
	} else {
		msg.Script = buffer
	}
	msg.MaxStdout = &maxstdout
	msg.MaxStderr = &maxstderr
	msg.Timeout = &timeout
	c.session.SendMessage("MasterService", "HandleRunScriptRequest", msg)
}
