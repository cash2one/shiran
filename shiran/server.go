package shiran

import (
	"net"
	"github.com/golang/glog"
	"reflect"
	"crypto/tls"
	"errors"
	"time"
	"runtime"
	"container/list"
)

var typeOfSession = reflect.TypeOf((*Session)(nil))

type methodType struct {
	method     reflect.Method
	ArgType    reflect.Type
}

type Service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

type Server struct {
	sessions	map[string]*Session
	register	chan *Session
	unregister	chan *Session
	services	map[string]*Service
	connectedCallbacks  *list.List
	aesKey		[]byte
}

func NewServer(aesKey []byte) *Server {
	return &Server{
		sessions:		make(map[string]*Session),
		register:		make(chan *Session),
		unregister:		make(chan *Session),
		services:		make(map[string]*Service),
		connectedCallbacks: list.New(),
		aesKey:			aesKey,
	}
}

func (server *Server) serviceGet(name string) *Service {
	if s, ok := server.services[name]; ok {
		return s
	}
	return nil
}

func (server *Server) RegisterService(rcvr interface{}) error {
	s := new(Service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if sname == "" {
		s := "RegisterService: no service name for type " + s.typ.String()
		glog.Info(s)
		return errors.New(s)
	}
	if _, present := server.services[sname]; present {
		return errors.New("RegisterService: service already defined: " + sname)
	}
	s.name = sname
	// Install the methods
	s.method = registerMethods(s.typ)
	server.services[s.name] = s
	return nil
}

// registerMethods returns suitable methods of typ
func registerMethods(typ reflect.Type) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}

		// Method needs two ins: receiver, *args, net.Conn.
		if mtype.NumIn() != 3 {
			glog.Warningf("registerMethods: method", mname, "has wrong number of ins:", mtype.NumIn())
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		// Second arg must be a pointer.
		sessionType := mtype.In(2)
		if sessionType.Kind() != reflect.Ptr {
			glog.Warningf("registerMethods: method", mname, "session type not a pointer:", sessionType)
			continue
		}
		if sessionType != typeOfSession {
			glog.Warningf("registerMethods: method", mname, "session type not Session:", sessionType, typeOfSession)
			continue
		}

		methods[mname] = &methodType{method: method, ArgType: argType}
	}
	return methods
}

func (server *Server) AddConnectedCallback(ccb ConnectedCallback) {
	server.connectedCallbacks.PushBack(ccb)
}


func (server *Server) handleSessions() {
	ticks := time.Tick(time.Second * 1)
	for {
		select {
		case session := <-server.register:
			server.sessions[session.Name] = session
		case session := <-server.unregister:
			delete(server.sessions, session.Name)
			close(session.packetQueue)
		case _ = <-ticks:
			glog.Infof("%d %d", len(server.sessions), runtime.NumGoroutine())
		}
	}
}

func (server *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Errorf("ListenAndServe: failed addr:%s %v", addr, err)
		return err
	}

	go server.serve(l)
	server.handleSessions()
	return nil
}

func (server *Server) TlsListenAndServe(addr string, tlsConfig *tls.Config) error {
	l, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		glog.Errorf("TlsListenAndServe: with tlsConfig failed addr:%s %v", addr, err)
		return err
	}

	go server.serve(l)
	server.handleSessions()
	return nil
}

func (server *Server) connected(session *Session) {
	for e := server.connectedCallbacks.Front(); e != nil; e = e.Next() {
		cb, ok := e.Value.(ConnectedCallback)
		if ok {
			cb(session)
		}
	}
}

func (server *Server) serve(l net.Listener) error {
	defer l.Close()

	var id int64
	var tempDelay time.Duration
	for {
		c, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				glog.Errorf("serve: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			glog.Errorf("serve: Accept failed %v", err)
			return err
		}
		tempDelay = 0
		id++
		go server.serveConn(c, id)
	}
}

func (server *Server) serveConn(conn net.Conn, id int64) {
	session := NewSession(id, conn, server, server.aesKey)
	server.register <- session
	defer func() { server.unregister <- session }()

	go session.sendPacketQueue()
	server.connected(session)
	session.recvMessage()
}

