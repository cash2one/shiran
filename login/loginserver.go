package login

import (
	"time"
	"runtime"
	"github.com/golang/glog"
	"github.com/williammuji/shiran/shiran"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Options struct {
	CertificateFile		string
	PrivateKeyFile		string
	CaFile				string
	ListenAddress		string
	ListenGateAddress	string
}

type LoginServer struct {
	//user
	register			chan *shiran.Session
	unregister			chan *shiran.Session
	server				*shiran.Server
	sessions			map[string]*shiran.Session
	//gate
	gateRegister        chan *shiran.Session
	gateUnregister      chan *shiran.Session
	gateServer          *shiran.Server
	gateSessions        map[string]*shiran.Session
	//service
	loginService		*LoginService
	gateService			*GateService
	//db
	db					*sql.DB
	//gatemanager
	gateManager			*GateManager
}

func NewLoginServer(opt *Options, db *sql.DB) *LoginServer {
	loginServer := &LoginServer{
		register:			make(chan *shiran.Session),
		unregister:			make(chan *shiran.Session),
		sessions:			make(map[string]*shiran.Session),
		gateRegister:		make(chan *shiran.Session),
		gateUnregister:		make(chan *shiran.Session),
		gateSessions:		make(map[string]*shiran.Session),
		db:					db,
		gateManager:		NewGateManager(),
	}
	loginServer.loginService = NewLoginService(db, loginServer.gateManager)
	loginServer.gateService = NewGateService(loginServer.gateManager)
	loginServer.server = shiran.NewServer(loginServer.register, loginServer.unregister, nil)
	loginServer.gateServer = shiran.NewServer(loginServer.gateRegister, loginServer.gateUnregister, nil)
	return loginServer 
}

func (loginServer *LoginServer) Run(opt *Options) {
	go loginServer.gateManager.eventLoop()

	loginServer.server.RegisterService(loginServer.loginService)
	loginServer.gateServer.RegisterService(loginServer.gateService)

	go loginServer.timer()

	go loginServer.gateServer.ListenAndServe(opt.ListenGateAddress)

	glog.Errorf("%s %s %s %s", opt.CertificateFile, opt.PrivateKeyFile, opt.CaFile, opt.ListenAddress)
	tlsConfig := shiran.GetServerTlsConfiguration(opt.CertificateFile, opt.PrivateKeyFile, opt.CaFile)
	loginServer.server.TlsListenAndServe(opt.ListenAddress, tlsConfig)

	//time.Sleep(30*time.Second);
}

func (loginServer *LoginServer) timer() {
	ticks := time.Tick(time.Second * 1)

	for {
		select {
		case session := <-loginServer.register:
			loginServer.sessions[session.Name] = session
		case session := <-loginServer.unregister:
			event := gateEvent{
				eventType:		ZONE_GATE_SESSION_UNREGISTER,
				session:		session,
			}
			loginServer.gateManager.PostEvent(event)
			delete(loginServer.sessions, session.Name)
		case session := <-loginServer.gateRegister:
			loginServer.gateSessions[session.Name] = session
		case session := <-loginServer.gateUnregister:
			delete(loginServer.gateSessions, session.Name)
		case _ = <-ticks:
			glog.Infof("%d %d %d", len(loginServer.sessions), len(loginServer.gateSessions), runtime.NumGoroutine())
		}
	}
}

