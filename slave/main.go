package main

import (
	"flag"
	"github.com/golang/glog"
	"os"
	"github.com/williammuji/shiran/masterslave"
)

type options struct {
	configFile  string
}

var opt options

func init() {
	flag.StringVar(&opt.configFile, "configFile", "./slaveconfig.json", "slave config json file")
}

func main() {
	if os.Getpid() == 0 || os.Geteuid() == 0 {
		glog.Info("Cannot run slave as root")
		return
	}

	flag.Parse()

	config := masterslave.NewSlaveConfig(opt.configFile)
	slave := masterslave.NewSlave(config)
	slave.Run()

	glog.Flush()
}

