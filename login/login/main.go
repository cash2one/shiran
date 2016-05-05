package main

import (
	"runtime"
	"os"
	"github.com/golang/glog"
	"flag"
	"runtime/pprof"
	"github.com/williammuji/shiran/login"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var opt login.Options

func init() {
	flag.StringVar(&opt.CertificateFile, "crt", "", "certificateFile")  //your/path/pki/issued/localhost.crt
	flag.StringVar(&opt.PrivateKeyFile, "key", "", "privateKeyFile")    //your/path/pki/private/localhost.key
	flag.StringVar(&opt.CaFile, "ca", "", "caFile")                     //your/path/pki/ca.crt
	flag.StringVar(&opt.ListenAddress, "listenAddress", "", "address listen for user")
	flag.StringVar(&opt.ListenGateAddress, "listenGateAddress", "", "address listen for gate")
}

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()
	defer glog.Flush()

	f, err := os.Create("cpu_loginserver.prof")
	if err != nil {
		glog.Errorf("ERROR os.Create cpu_loginserver.prof failed")
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	db, err := sql.Open("mysql", "noahsark:noahsark@tcp(127.0.0.1:3306)/shiran")
	if err != nil {
		glog.Errorf("sql.Open failed err:%s", err)
	}
	defer db.Close()

	loginServer := login.NewLoginServer(&opt, db)
	loginServer.Run(&opt)

}

