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
	flag.StringVar(&opt.configFile, "configFile", "./masterconfig.json", "master config json file")
}

func main() {
	if os.Getpid() == 0 || os.Geteuid() == 0 {
		glog.Info("Cannot run master as root")
		return
	}

	flag.Parse()

	config := masterslave.NewMasterConfig(opt.configFile)
	glog.Infof("start master")
	master := masterslave.NewMaster(config)
	master.Run()

	glog.Flush()
}

