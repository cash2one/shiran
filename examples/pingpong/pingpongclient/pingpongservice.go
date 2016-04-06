package pingpongclient

import (
	"github.com/williammuji/shiran/shiran"
	"github.com/williammuji/shiran/examples/pingpong"
	_ "github.com/golang/glog"
	"bytes"
	"errors"
	"sync/atomic"
)

type PingpongService struct {
	size			int
	Data			[]byte
	TotalMsgCount   int64
}

func NewPingpongService(size int) *PingpongService {
	pps := &PingpongService{
		size:       size,
		Data:		make([]byte, size),
	}
	for i := 0; i < pps.size; i++ {
		pps.Data[i] = "0123456789ABCDEF"[i%16]
	}
	return pps
}

func (pps *PingpongService) HandlePingPongData(msg *protocol.PingPongData, session *shiran.Session) {
	//glog.Infov("%v", msg)

	if !bytes.Equal(pps.Data, msg.Data) {
		panic(errors.New("Recv PingPongData Wrong"))
	}
	atomic.AddInt64(&pps.TotalMsgCount, int64(1))
	session.SendMessage("PingpongService", "HandlePingPongData", msg)
}
