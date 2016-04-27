package main

import (
	"flag"
	"github.com/golang/glog"
	"os"
	"strings"
	"github.com/williammuji/shiran/masterslave"
)

var opt masterslave.CommandOptions 

func init() {
	flag.StringVar(&opt.ConfigFile, "configFile", "./commanderconfig.json", "slave config json file")
	flag.StringVar(&opt.SlaveName, "slaveName", "", "slave Name")
	flag.BoolVar(&opt.Add, "add", false, "add")
	flag.BoolVar(&opt.Remove, "remove", false, "remove")
	flag.BoolVar(&opt.Start, "start", false, "start")
	flag.BoolVar(&opt.Stop, "stop", false, "stop")
	flag.BoolVar(&opt.Restart, "restart", false, "restart")
	flag.BoolVar(&opt.Get, "get", false, "get")
	flag.BoolVar(&opt.List, "list", false, "list")
	flag.StringVar(&opt.AppName, "appName", "", "app name")
}

func main() {
	if os.Getpid() == 0 || os.Geteuid() == 0 {
		glog.Info("Cannot run slave as root")
		return
	}

	flag.Parse()
	opt.AppNames = strings.Split(opt.AppName, ";")

	commander := masterslave.NewCommander(&opt)
	commander.Run()

	glog.Flush()
}

