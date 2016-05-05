package main

import (
	"runtime"
	"os"
	"github.com/golang/glog"
	"flag"
	"runtime/pprof"
	"github.com/williammuji/shiran/gate"
)

func init() {
	flag.StringVar(&gate.Opt.ListenAddress, "listenAddress", "", "address listen for user")
	flag.IntVar(&gate.Opt.Zone, "zone", 0, "zone id")
	flag.StringVar(&gate.Opt.LoginAddress, "loginAddress", "", "connect login address")
}

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()

	f, err := os.Create("cpu_gateserver.prof")
	if err != nil {
		glog.Errorf("ERROR os.Create cpu_gateserver.prof failed")
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	gateServer := gate.NewGateServer()
	gateServer.Run()

	glog.Flush()
}

