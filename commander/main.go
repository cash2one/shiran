package main

import (
	"flag"
	"github.com/golang/glog"
	"os"
	"strings"
	"github.com/williammuji/shiran/masterslave"
)

var opt masterslave.CommandOptions 
var files string
var args string

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
	flag.BoolVar(&opt.GetHardware, "getHardware", false, "get hardware")
	flag.BoolVar(&opt.Lshw, "lshw", false, "lshw")
	flag.BoolVar(&opt.GetFileContent, "getFileContent", false, "get file content")
	flag.StringVar(&opt.FileName, "fileName", "", "file name")
	flag.Int64Var(&opt.MaxSize, "maxSize", 0, "max file size")
	flag.BoolVar(&opt.GetFileChecksum, "getFileChecksum", false, "get file checksum")
	flag.StringVar(&files, "files", "", "files name")
	flag.BoolVar(&opt.RunCommand, "runCommand", false, "run command")
	flag.StringVar(&opt.Command, "command", "", "command")
	flag.StringVar(&args, "args", "", "args")
	flag.IntVar(&opt.MaxStdout, "maxStdout", 0, "max stdout size")
	flag.IntVar(&opt.MaxStderr, "maxStderr", 0, "max stderr size")
	flag.IntVar(&opt.Timeout, "timeout", 0, "timeout")
	flag.BoolVar(&opt.RunScript, "runScript", false, "run script")
	flag.StringVar(&opt.Script, "script", "", "script name")
}

func main() {
	if os.Getpid() == 0 || os.Geteuid() == 0 {
		glog.Info("Cannot run slave as root")
		return
	}

	flag.Parse()
	if opt.AppName != "" {
		opt.AppNames = strings.Split(opt.AppName, ";")
	}
	if files != "" {
		opt.Files = strings.Split(files, ";")
	}
	if args != "" {
		opt.Args = strings.Split(args, ";")
	}

	commander := masterslave.NewCommander(&opt)
	commander.Run()

	glog.Flush()
}

