package main

import (
	"time"
	"runtime"
	"os"
	"github.com/golang/glog"
	"flag"
	"runtime/pprof"
	"github.com/williammuji/shiran/shiran"
	"github.com/williammuji/shiran/examples/pingpong"
	"github.com/williammuji/shiran/examples/pingpong/pingpongclient"
	"sync/atomic"
)

type options struct {
	tls			bool
	aesKey		string
	num			int64
	size		int
	timeout		int
	caFile      string
}

var opt options

func init() {
	flag.BoolVar(&opt.tls, "tls", false, "tls conn")
	flag.StringVar(&opt.aesKey, "aesKey", "", "AES key")
	flag.Int64Var(&opt.num, "n", 5000, "launch client num")
	flag.IntVar(&opt.size, "s", 4096, "client send msg size")
	flag.IntVar(&opt.timeout, "t", 20, "client quit timeout")
	flag.StringVar(&opt.caFile, "ca", "", "caFile")	//your/path/pki/ca.crt
}

type PingpongClient struct {
	register		chan *shiran.Session
	unregister		chan *shiran.Session
	client			*shiran.Client
	ppService		*pingpongclient.PingpongService
	sessions		map[string]*shiran.Session
	sessionCount	int64
}

func NewPingpongClient(opt *options) *PingpongClient {
	ppClient := &PingpongClient{
		register:   make(chan *shiran.Session),
		unregister: make(chan *shiran.Session),
		ppService:	pingpongclient.NewPingpongService(opt.size),
		sessions:   make(map[string]*shiran.Session),
	}
	ppClient.client = shiran.NewClient("localhost:8848", opt.num, ppClient.register, ppClient.unregister, []byte(opt.aesKey))
	return ppClient
}

func (ppClient *PingpongClient) run(opt *options) {
	ppClient.client.RegisterService(ppClient.ppService)
	
	if opt.caFile != "" {
		tlsConfig := shiran.GetClientTlsConfiguration(opt.caFile)
		go ppClient.client.TlsConnectServer(tlsConfig)
	} else {
		go ppClient.client.ConnectServer()
	}

	ppClient.timer(opt)
}

func (ppClient *PingpongClient) timer(opt *options) {
	ticks := time.Tick(time.Second * 1)
	quitTicks := time.NewTimer(time.Duration(opt.timeout) * time.Second)
	var quit bool

	for {
		select {
		case session := <-ppClient.register:
			ppClient.sessions[session.Name] = session
			ppClient.onConnection(session)
		case session := <-ppClient.unregister:
			delete(ppClient.sessions, session.Name)
			ppClient.onConnection(session)
		case _ = <-ticks:
			glog.Infof("%d %d", len(ppClient.sessions), runtime.NumGoroutine())
			if quit == true && len(ppClient.sessions) == 0 {
				return
			}
		case _ = <-quitTicks.C:
			for _, session := range ppClient.sessions {
				session.Close()
			}
			quit = true
		}
	}
}

func (ppClient *PingpongClient) onConnection(session *shiran.Session) {
	if session.Closed == false {
		msg := &protocol.PingPongData{Data: ppClient.ppService.Data}
		//glog.Infof("%v", msg)
		session.SendMessage("PingpongService", "HandlePingPongData", msg)
		ppClient.sessionCount++
		if ppClient.sessionCount == opt.num {
			glog.Infof("All %d sessions connected", opt.num)
		}
	} else if session.Closed == true {
		ppClient.sessionCount--
		if ppClient.sessionCount == 0 {
			glog.Infof("All %d sessions disconnected", opt.num)
			total := atomic.LoadInt64(&ppClient.ppService.TotalMsgCount)
			glog.Infof("timer: PingpongService timeout throughput: %.3f MiB/s", float64(total * int64(len(ppClient.ppService.Data)))/(1024.0*1024.0)/float64(opt.timeout))
		}
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()

	f, err := os.Create("cpu_pingpongclient.prof")
	if err != nil {
		glog.Errorf("ERROR os.Create cpu_pingpongclient.prof failed")
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	ppClient := NewPingpongClient(&opt)
	ppClient.run(&opt)

	glog.Flush()
}
