package pingpongclient

import (
	"github.com/williammuji/shiran2/shiran"
	"github.com/williammuji/shiran2/examples/pingpong"
	"github.com/golang/glog"
	"bytes"
	"errors"
	"sync/atomic"
	"time"
)

type PingPongService struct {
	client			*shiran.Client
	size			int
	timeout			int
	data			[]byte
	totalMsgCount   int64
}

func NewPingPongService(client *shiran.Client, size, timeout int) *PingPongService {
	pps := &PingPongService{
		client:     client,
		size:       size,
		timeout:    timeout,
		data:        make([]byte, size),
	}
	for i := 0; i < pps.size; i++ {
		pps.data[i] = "0123456789ABCDEF"[i%16]
	}
	client.AddConnectedCallback(pps.SessionConnected)
	go pps.timer()
	return pps
}

func (pps *PingPongService) HandlePingPongData(msg *protocol.PingPongData, session *shiran.Session) {
	//glog.Infov("%v", msg)

	if !bytes.Equal(pps.data, msg.Data) {
		panic(errors.New("Recv PingPongData Wrong"))
	}
	atomic.AddInt64(&pps.totalMsgCount, int64(1))
	session.SendMessage("PingPongService", "HandlePingPongData", msg)
}

func (pps *PingPongService) timer() {
	t := time.NewTimer(time.Duration(pps.timeout) * time.Second)

	for {
		select {
		case <-t.C:
			total := atomic.LoadInt64(&pps.totalMsgCount)
			glog.Infof("timer: PingPongService timeout throughput: %.3f MiB/s", float64(total * int64(len(pps.data)))/(1024.0*1024.0)/float64(pps.timeout))
			pps.client.Quit()
			break
		}
	}
}

func (pps *PingPongService) SessionConnected(session *shiran.Session) {
	msg := &protocol.PingPongData{Data: pps.data}
	//glog.Infof("%v", msg)
	session.SendMessage("PingPongService", "HandlePingPongData", msg)
}

