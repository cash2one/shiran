package login

import (
	"github.com/williammuji/shiran/shiran"
	. "github.com/williammuji/shiran/proto/userproto"
	. "github.com/williammuji/shiran/proto/gatelogin"
	"github.com/golang/glog"
	"crypto/sha256"
	"math/rand"
	"database/sql"
	"bytes"
	"time"
)

type LoginService struct {
	db				*sql.DB
	gateManager		*GateManager
}

func NewLoginService(db *sql.DB, gateManager *GateManager) *LoginService {
	ls := &LoginService{
		db:				db,
		gateManager:	gateManager,
	}
	return ls
}

func (ls *LoginService) HandleUserLoginLoginRequest(request *UserLoginLoginRequest, session *shiran.Session) {
	glog.Infof("%+v", request)

	response := &UserLoginLoginResponse{}
	stmt, err := ls.db.Prepare("select PASSWD, SALT from ACCOUNT where NAME = ?")
	if err != nil {
		glog.Errorf("%v", err)
		response.State = UserLoginLoginState_kDatabaseError.Enum()
		session.SendMessage("LoginService", "HandleUserLoginLoginResponse", response)
		return
	}

	var passwd, salt string
	err = stmt.QueryRow(request.GetName()).Scan(&passwd, &salt)
	if err != nil {
		glog.Errorf("%v", err)
		response.State = UserLoginLoginState_kNamePasswdError.Enum()
		session.SendMessage("LoginService", "HandleUserLoginLoginResponse", response)
		return
	}

	str := request.GetPasswd() + salt
	res := sha256.Sum256([]byte(str))
	hashed := res[:]
	if bytes.Equal(hashed, []byte(passwd)) {
		response.State = UserLoginLoginState_kNamePasswdError.Enum()
		session.SendMessage("LoginService", "HandleUserLoginLoginResponse", response)
		return
	}


	//forward to gateserver
	rkRequest := &UserRandomKeyRequest{
		Zone:		request.Zone,
		Name:		request.Name,
	}
	randomKey := make([]byte, 16)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		randomKey[i] = byte(r.Intn(256))
	}
	rkRequest.RandomKey = randomKey

	event := gateEvent{
		eventType:		USER_LOGIN_REQUEST,
		msg:			rkRequest,
		session:		session,
	}
	ls.gateManager.PostEvent(event)
}
