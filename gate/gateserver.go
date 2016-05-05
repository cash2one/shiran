package gate

import (
	"time"
	"runtime"
	"github.com/golang/glog"
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/gatelogin"
)

type Options struct {
	ListenAddress		string
	Zone				int
	LoginAddress		string
}

var Opt Options

type GateServer struct {
	//user
	register			chan *shiran.Session
	unregister			chan *shiran.Session
	server				*shiran.Server
	sessions			map[string]*shiran.Session
	//login
	loginRegister		chan *shiran.Session
	loginUnregister     chan *shiran.Session
	loginClient         *shiran.Client
	loginSessions       map[string]*shiran.Session
	//service
	gateService			*GateService
	//userManager
	userManager			*UserManager
}

func NewGateServer()  *GateServer {
	gateServer := &GateServer{
		register:			make(chan *shiran.Session),
		unregister:			make(chan *shiran.Session),
		sessions:			make(map[string]*shiran.Session),
		loginRegister:		make(chan *shiran.Session),
		loginUnregister:	make(chan *shiran.Session),
		loginSessions:		make(map[string]*shiran.Session),
		userManager:		NewUserManager(),
	}
	gateServer.gateService = NewGateService(gateServer.userManager)
	gateServer.server = shiran.NewServer(gateServer.register, gateServer.unregister, nil)
	gateServer.loginClient = shiran.NewClient(Opt.LoginAddress, 1, gateServer.loginRegister, gateServer.loginUnregister, nil)
	return gateServer 
}

func (gateServer *GateServer) Run() {
	gateServer.server.RegisterService(gateServer.gateService)
	gateServer.loginClient.RegisterService(gateServer.gateService)

	go gateServer.timer()

	go gateServer.loginClient.ConnectServer()

	gateServer.server.ListenAndServe(Opt.ListenAddress)
}

func (gateServer *GateServer) timer() {
	ticks := time.Tick(time.Second * 1)

	for {
		select {
		case session := <-gateServer.register:
			gateServer.sessions[session.Name] = session
		case session := <-gateServer.unregister:
			event := userEvent{
				eventType:      USER_UNREGISTER,
				session:        session,
			}
			gateServer.userManager.PostEvent(event)
			delete(gateServer.sessions, session.Name)
		case session := <-gateServer.loginRegister:
			gateServer.loginSessions[session.Name] = session
			gateServer.onLoginConnection(session)
		case session := <-gateServer.loginUnregister:
			delete(gateServer.loginSessions, session.Name)
		case _ = <-ticks:
			glog.Infof("%d %d %d", len(gateServer.sessions), len(gateServer.loginSessions), runtime.NumGoroutine())
		}
	}
}

func (gateServer *GateServer) onLoginConnection(session *shiran.Session) {
	request := &GateLoginRequest{
		GateListenAddress:		&Opt.ListenAddress,
	}
	var zoneID int32 = int32(Opt.Zone)
	request.Zone = &zoneID

	glog.Errorf("%v %d %v", Opt.ListenAddress, Opt.Zone, request)

	session.SendMessage("GateService", "HandleGateLoginRequest", request)
}
