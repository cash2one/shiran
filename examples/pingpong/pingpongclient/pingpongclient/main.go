package main

import (
	"runtime"
	"os"
	"github.com/golang/glog"
	"flag"
	"runtime/pprof"
	"github.com/williammuji/shiran2/shiran"
	"github.com/williammuji/shiran2/examples/pingpong/pingpongclient"
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

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()
	//logFile, err := os.Create("pingpongclient.log")
	//defer logFile.Close()
	//if err != nil {
	//	log.Fatalln("open log file error")
	//}
	//log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	//log.SetOutput(logFile)

	f, err := os.Create("cpu_pingpongclient.prof")
	if err != nil {
		glog.Errorf("ERROR os.Create cpu_pingpongclient.prof failed")
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	client := shiran.NewClient("localhost:8848", opt.num, []byte(opt.aesKey))
	ppservice := pingpongclient.NewPingPongService(client, opt.size, opt.timeout)
	client.RegisterService(ppservice)
	if opt.caFile != "" {
		tlsConfig := shiran.GetClientTlsConfiguration(opt.caFile)
		client.TlsConnectServer(tlsConfig)
	} else {
		client.ConnectServer()
	}

	glog.Flush()
}
