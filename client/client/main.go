package main

import (
	"flag"
	"github.com/golang/glog"
	"os"
	"github.com/williammuji/shiran/client"
	"time"
	"runtime"
)

type options struct {
	configFile  string
}

var opt options

func init() {
	flag.StringVar(&opt.configFile, "configFile", "./clientconfig.json", "client config json file")
}

func main() {
	if os.Getpid() == 0 || os.Geteuid() == 0 {
		glog.Error("Cannot run client as root")
		return
	}

	flag.Parse()

	config := client.NewClientConfig(opt.configFile)
	if len(config.LoginAddress) == 0 || len(config.Accounts) == 0 {
		glog.Error("Cannot run client if null config")
		return
	}

	index := 0
	for  i, _ := range config.Accounts {
		loginClient := client.NewLoginClient(config.LoginAddress[index], &config.Accounts[i])
		index++
		if index == len(config.LoginAddress) {
			index = 0
		}
		go loginClient.Run(config.CaFile)
	}

	ticks := time.Tick(time.Second * 1)
	for {
		select {
		case _ = <-ticks:
			glog.Infof("%d", runtime.NumGoroutine())
		}
	}

	glog.Flush()
}

