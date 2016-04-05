package main

import (
	"runtime"
	"os"
	"github.com/golang/glog"
	"flag"
	"errors"
	"runtime/pprof"
	"github.com/williammuji/shiran2/shiran"
	"github.com/williammuji/shiran2/examples/pingpong/pingpongserver"
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

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()
	//logFile, err := os.Create("pingpongserver.log")
	//defer logFile.Close()
	//if err != nil {
	//	log.Fatalln("open log file error")
	//}
	//log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	//log.SetOutput(logFile)

	f, err := os.Create("cpu_pingpongserver.prof")
	if err != nil {
		glog.Errorf("ERROR os.Create cpu_pingpongserver.prof failed")
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	server := shiran.NewServer([]byte(opt.aesKey))
	server.RegisterService(pingpongserver.NewPingPongService())
	if opt.caFile != "" {
		if opt.privateKeyFile == "" || opt.certificateFile == "" {
			panic(errors.New("key crt null"))
		}
		tlsConfig := shiran.GetServerTlsConfiguration(opt.certificateFile, opt.privateKeyFile, opt.caFile)
		server.TlsListenAndServe("localhost:8848", tlsConfig)
	} else {
		server.ListenAndServe("localhost:8848")
	}

	glog.Flush()
}
