package shiran 

import (
	"github.com/golang/glog"
	"net"
	"runtime"
	"time"
	"reflect"
	"errors"
	"strings"
	"crypto/tls"
	"container/list"
)

type Client struct {
	ipport			string
	num				int64
	sessions		map[string]*Session
	register		chan *Session
	unregister		chan *Session
	services		map[string]*Service
	connectedCallbacks  *list.List
	aesKey			[]byte
	quit			chan bool
}

func NewClient(ipport string, num int64, aesKey []byte) *Client {
	return &Client{
		ipport:		ipport,
		num:		num,
		sessions:   make(map[string]*Session),
		register:   make(chan *Session),
		unregister: make(chan *Session),
		services:   make(map[string]*Service),
		connectedCallbacks: list.New(),
		aesKey:		aesKey,
		quit:		make(chan bool),
	}
}

func (client *Client) AddConnectedCallback(ccb ConnectedCallback) {
	client.connectedCallbacks.PushBack(ccb)
}

func (client *Client) serviceGet(name string) *Service {
	if s, ok := client.services[name]; ok {
		return s
	}
	return nil
}

func (client *Client) RegisterService(rcvr interface{}) error {
	s := new(Service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if sname == "" {
		s := "RegisterService: no service name for type " + s.typ.String()
		glog.Info(s)
		return errors.New(s)
	}
	if _, present := client.services[sname]; present {
		return errors.New("RegisterService: service already defined: " + sname)
	}
	s.name = sname
	// Install the methods
	s.method = registerMethods(s.typ)
	client.services[s.name] = s
	return nil
}

func (client *Client) handleSessions() {
	ticks := time.Tick(time.Second * 1)
	var quit bool
	for {
		select {
		case session := <-client.register:
			client.sessions[session.Name] = session
		case session := <-client.unregister:
			delete(client.sessions, session.Name)
			close(session.packetQueue)
		case _ = <-ticks:
			glog.Infof("%d %d", len(client.sessions), runtime.NumGoroutine())
			if quit == true && len(client.sessions) == 0 {
				return
			}
		case _ = <-client.quit:
			for _, session := range client.sessions {
				session.close()
			}
			quit = true
		}
	}
}

func (client *Client) ConnectServer() {
	if client.num > 1 {
		for i := int64(1); i <= client.num; i++ {
			conn, err := net.Dial("tcp", client.ipport)
			if err != nil {
				glog.Errorf("ConnectServer: failed ipport:%s id:%d err:%v", client.ipport, i, err)
				continue
			}
			go client.serveConn(conn, i)
		}
	} else {
		ipports := strings.Split(client.ipport, ";")
		if len(ipports) == 0 || ipports[0] == "" {
			glog.Errorf("ConnectServer: failed ipport:%s", client.ipport)
		}
		for i, addr := range ipports {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				glog.Errorf("ConnectServer: failed ipport:%s id:%d err:%v", addr, i+1, err)
				continue
			}
			go client.serveConn(conn, int64(i+1))
		}
	}

	client.handleSessions()
}

func (client *Client) TlsConnectServer(tlsConfig *tls.Config) {
	if client.num > 1 {
		for i := int64(1); i <= client.num; i++ {
			conn, err := tls.Dial("tcp", client.ipport, tlsConfig)
			if err != nil {
				glog.Errorf("TlsConnectServer: failed ipport:%s id:%d err:%v", client.ipport, i, err)
				continue
			}
			err = conn.Handshake()
			if err != nil {
				glog.Errorf("TlsConnectServer: handshake failed addr:%s id:%d err:%v", client.ipport, i, err)
				continue    
			}
			go client.serveConn(conn, i)
		}
	} else {
		ipports := strings.Split(client.ipport, ";")
		if len(ipports) == 0 || ipports[0] == "" {
			glog.Errorf("TlsConnectServer: failed ipport:%s", client.ipport)
		}
		for i, addr := range ipports {
			conn, err := tls.Dial("tcp", addr, tlsConfig)
			if err != nil {
				glog.Errorf("TlsConnectServer: failed ipport:%s id:%d err:%v", addr, i+1, err)
				continue
			}
			err = conn.Handshake()
			if err != nil {
				glog.Errorf("TlsConnectServer: handshake failed addr:%s id:%d err:%v", addr, i, err)
				continue    
			}
			go client.serveConn(conn, int64(i+1))
		}
	}

	client.handleSessions()
}

func (client *Client) connected(session *Session) {
	for e := client.connectedCallbacks.Front(); e != nil; e = e.Next() {
		cb, ok := e.Value.(ConnectedCallback)
		if ok {
			cb(session)
		}
	}
}

func (client *Client) serveConn(conn net.Conn, id int64) {
	session := NewSession(id, conn, client, client.aesKey)
	client.register <- session
	defer func() { client.unregister <- session }()

	go session.sendPacketQueue()
	client.connected(session)
	session.recvMessage()
}

func (client *Client) Quit() {
	client.quit <- true
}
