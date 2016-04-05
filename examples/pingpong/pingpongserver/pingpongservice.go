package pingpongserver

import (
	"github.com/williammuji/shiran2/shiran"
	"github.com/williammuji/shiran2/examples/pingpong"
	_ "log"
)

type PingPongService struct {
}

func NewPingPongService() *PingPongService {
	pps := &PingPongService{
	}
	return pps
}

func (pps *PingPongService) HandlePingPongData(msg *protocol.PingPongData, session *shiran.Session) {
	//glog.Infof("%v", msg)

	session.SendMessage("PingPongService", "HandlePingPongData", msg)
}

