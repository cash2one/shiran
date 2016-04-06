package pingpongserver

import (
	"github.com/williammuji/shiran/shiran"
	"github.com/williammuji/shiran/examples/pingpong"
	_ "log"
)

type PingpongService struct {
}

func NewPingpongService() *PingpongService {
	pps := &PingpongService{
	}
	return pps
}

func (pps *PingpongService) HandlePingPongData(msg *protocol.PingPongData, session *shiran.Session) {
	//glog.Infof("%v", msg)

	session.SendMessage("PingpongService", "HandlePingPongData", msg)
}

