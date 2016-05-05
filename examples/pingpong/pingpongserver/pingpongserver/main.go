package main

import (
	"time"
	"runtime"
	"os"
	"github.com/golang/glog"
	"flag"
	"errors"
	"runtime/pprof"
	"github.com/williammuji/shiran/shiran"
	"github.com/williammuji/shiran/examples/pingpong/pingpongserver"
)

type options struct {
	aesKey			string
	certificateFile string
	privateKeyFile  string
	caFile          string
}

var opt options

func init() {
	flag.StringVar(&opt.aesKey, "aesKey", "", "AES key")
	flag.StringVar(&opt.certificateFile, "crt", "", "certificateFile")	//your/path/pki/issued/localhost.crt
	flag.StringVar(&opt.privateKeyFile, "key", "", "privateKeyFile")	//your/path/pki/private/localhost.key
	flag.StringVar(&opt.caFile, "ca", "", "caFile")						//your/path/pki/ca.crt
}

type PingpongServer struct {
	register    chan *shiran.Session
	unregister  chan *shiran.Session
	server		*shiran.Server
	ppService   *pingpongserver.PingpongService
	sessions    map[string]*shiran.Session
}

func NewPingpongServer(opt *options) *PingpongServer {
	ppServer := &PingpongServer{
		register:   make(chan *shiran.Session),
		unregister: make(chan *shiran.Session),
		ppService:  pingpongserver.NewPingpongService(),
		sessions:   make(map[string]*shiran.Session),
	}
	ppServer.server = shiran.NewServer(ppServer.register, ppServer.unregister, []byte(opt.aesKey))
	return ppServer 
}

func (ppServer *PingpongServer) run(opt *options) {
	ppServer.server.RegisterService(ppServer.ppService)

	go ppServer.timer()

	if opt.caFile != "" {
		if opt.privateKeyFile == "" || opt.certificateFile == "" {
			panic(errors.New("key crt null"))
		}
		tlsConfig := shiran.GetServerTlsConfiguration(opt.certificateFile, opt.privateKeyFile, opt.caFile)
		ppServer.server.TlsListenAndServe("127.0.0.1:8848", tlsConfig)
	} else {
		ppServer.server.ListenAndServe("127.0.0.1:8848")
	}
}

func (ppServer *PingpongServer) timer() {
	ticks := time.Tick(time.Second * 1)

	for {
		select {
		case session := <-ppServer.register:
			ppServer.sessions[session.Name] = session
		case session := <-ppServer.unregister:
			delete(ppServer.sessions, session.Name)
		case _ = <-ticks:
			glog.Infof("%d %d", len(ppServer.sessions), runtime.NumGoroutine())
		}
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()

	f, err := os.Create("cpu_pingpongserver.prof")
	if err != nil {
		glog.Errorf("ERROR os.Create cpu_pingpongserver.prof failed")
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	ppServer := NewPingpongServer(&opt)
	ppServer.run(&opt)

	glog.Flush()
}
