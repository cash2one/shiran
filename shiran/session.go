package shiran 

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/glog"
	"net"
	"strconv"
	"reflect"
)

type ConnectedCallback func(*Session)

type ServiceGetter interface {
	serviceGet(name string) *Service 
}

type Session struct {
	Name        string
	conn        net.Conn
	packetQueue chan []byte
	codec		Codec
	serviceGetter	ServiceGetter
	closed		bool
}

func NewSession(id int64, conn net.Conn, serviceGetter ServiceGetter, aesKey []byte) *Session {
	session := &Session{
		Name:        conn.RemoteAddr().String() + "_ID_" + strconv.FormatInt(id, 10),
		conn:        conn,
		packetQueue: make(chan []byte, 1024),	//size?
		serviceGetter:	serviceGetter,
	}
	if len(aesKey) > 0 {
		session.codec = NewSessionAesCodec(session.conn, aesKey)
	} else {
		session.codec = NewSessionCodec(session.conn)
	}
	glog.Infof("Session.UP   Name:%s LocalAddr:%s RemoteAddr:%s", session.Name, session.conn.LocalAddr().String(), session.conn.RemoteAddr().String())
	return session
}

func (session *Session) close() {
	if session.closed == false {
		glog.Infof("Session.DOWN Name:%s LocalAddr:%s RemoteAddr:%s", session.Name, session.conn.LocalAddr().String(), session.conn.RemoteAddr().String())
		session.conn.Close()
		session.closed = true
	}
}

func (session *Session) sendPacketQueue() {
	defer session.close()

	for packet := range session.packetQueue {
		n, err := session.conn.Write(packet)
		if err != nil {
			glog.Errorf("sendPacketQueue: session.Name:%s %v", session.Name, err)
			break
		}
		if n != len(packet) {
			glog.Errorf("sendPacketQueue: session.Name:%s n:%d != len(packet):%d", session.Name, n, len(packet))
			break
		}
	}
}

func (session *Session) SendMessage(service, method string, pb proto.Message) error {
	shiranMsg := new(ShiranMessage)
	shiranMsg.Service = &service
	shiranMsg.Method = &method

	msg, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	shiranMsg.Msg = msg
	payload, err := proto.Marshal(shiranMsg)
	if err != nil {
		return err
	}

	packet, err := session.codec.Encode(payload)
	if err != nil {
		return err
	}

	session.packetQueue <- packet
	return nil
}

func (session *Session) recvMessage() {
	for {
		payload, err := session.codec.Decode()
		if err != nil {
			glog.Errorf("recvMessage: codec.Decode ShiranMessage err: %v", err)
			break	
		}

		shiranMsg := new(ShiranMessage)
		err = proto.Unmarshal(payload, shiranMsg)
		if err != nil {
			glog.Errorf("recvMessage: proto.Unmarshal ShiranMessage err: %v", err)
			break	
		}

		service := session.serviceGetter.serviceGet(shiranMsg.GetService())
		if service == nil {
			glog.Errorf("recvMessage: can't find service %s", shiranMsg.GetService())
			break	
		}
		mtype := service.method[shiranMsg.GetMethod()]
		if mtype == nil {
			glog.Errorf("recvMessage: can't find method %s", shiranMsg.GetMethod())
			break	
		}

		// Decode the argument value.
		var argv reflect.Value
		argIsValue := false // if true, need to indirect before calling.
		if mtype.ArgType.Kind() == reflect.Ptr {
			argv = reflect.New(mtype.ArgType.Elem())
		} else {
			argv = reflect.New(mtype.ArgType)
			argIsValue = true
		}

		msg, ok := argv.Interface().(proto.Message)
		if !ok {
			glog.Errorf("recvMessage: body is not a protobuf Message")
			break	
		}
		err = proto.Unmarshal(shiranMsg.GetMsg(), msg)
		if err != nil {
			glog.Errorf("recvMessage: proto.Unmarshal msg failed err: %v", err)
			break	
		}
		if argIsValue {
			argv = argv.Elem()
		}

		function := mtype.method.Func
		// Invoke the method, providing a new value for the reply.
		function.Call([]reflect.Value{service.rcvr, argv, reflect.ValueOf(session)})
	}
}
