package masterslave

import (
	"os"
	"time"
	"io/ioutil"
	"github.com/golang/glog"
	"github.com/williammuji/shiran/shiran"
	"syscall"
)

type Heartbeat struct {
	name			string
	startTime		int64
	interval		int
	session			*shiran.Session
	slaveService	*SlaveService
	quit			chan bool
}

func NewHeartbeat(slaveConfig *SlaveConfig, session *shiran.Session, slaveService *SlaveService) *Heartbeat {
	hb := &Heartbeat{
		name:			slaveConfig.Name,
		startTime:		int64(time.Now().UnixNano()/1000),
		interval:		slaveConfig.HeartbeatInterval,
		session:		session,
		slaveService:	slaveService,
		quit:			make(chan bool),
	}
	return hb
}

func (hb *Heartbeat) start() {
	event := Event{
		eventType:      HEARTBEAT,
		msg:            hb.beat(true),
		session:        hb.session,
	}
	hb.slaveService.appManager.PostEvent(event)

	go hb.timer()
}

func (hb *Heartbeat) stop() {
	hb.quit <- true
}

func (hb *Heartbeat) timer() {
	ticks := time.Tick(time.Second * time.Duration(hb.interval))

	for {
		select {
		case _ = <-ticks:
			event := Event{
				eventType:      HEARTBEAT,
				msg:            hb.beat(false),
				session:        hb.session,
			}
			hb.slaveService.appManager.PostEvent(event)
		case _ = <-hb.quit:
			return
		}
	}
}

func (hb *Heartbeat) beat(showStatic bool) *SlaveHeartbeat {
	msg := &SlaveHeartbeat{
		SlaveName:			&hb.name,
		StartTimeUs:		&hb.startTime,
	}
	var cur int64 = int64(time.Now().UnixNano()/1000)
	msg.SendTimeUs = &cur

	if showStatic == true {
		name, err := os.Hostname()
		if err == nil {
			msg.HostName = &name
		}

		var pid int32 = int32(os.Getpid())
		msg.SlavePid = &pid

		msg.Cpuinfo = readFile("/proc/cpuinfo")
		msg.Version = readFile("/proc/version")
		msg.EtcMtab = readFile("/etc/mtab")

		un := syscall.Utsname{}
		err = syscall.Uname(&un)
		if err == nil {
			msg.Uname = &SlaveHeartbeat_Uname{
				Sysname:		B2S(un.Sysname[:]),
				Nodename:		B2S(un.Nodename[:]),
				Release:		B2S(un.Release[:]),
				Version:		B2S(un.Version[:]),
				Machine:		B2S(un.Machine[:]),
				Domainname:		B2S(un.Domainname[:]),
			}
		}
	}

	msg.Meminfo = readFile("/proc/meminfo")
	msg.ProcStat = readFile("/proc/stat")
	msg.Loadavg = readFile("/proc/loadavg")
	msg.Diskstats = readFile("/proc/diskstats")
	msg.NetDev = readFile("/proc/net/dev")
	msg.NetTcp = readFile("/proc/net/tcp")

	return msg
}

func readFile(path string) *string {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Errorf("readFile %s failed err:%s", path, err)
		return nil
	} else {
		var res string
		res = string(buffer)
		return &res
	}
}

func B2S(bs []int8) *string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}
	res := string(b)
	return &res
}
