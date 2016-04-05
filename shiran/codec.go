package shiran 

import (
	"io"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/glog"
)

type Codec interface {
	Decode() ([]byte, error)
	Encode(payload []byte) ([]byte, error)
}

type SessionCodec struct {
	r       io.Reader
}

func NewSessionCodec(r io.Reader) *SessionCodec {
	return &SessionCodec{
		r:	r,
	}
}

func (codec *SessionCodec) Decode() (payload []byte, err error) {
	header := make([]byte, 4)
	_, err = io.ReadFull(codec.r, header)
	if err != nil {
		return
	}

	length := binary.BigEndian.Uint32(header)
	if length > 64*1024*1024 {
		fmt.Errorf("Decode: invalid length %d", length)
		return
	}

	payload = make([]byte, length)
	_, err = io.ReadFull(codec.r, payload)
	if err != nil {
		return
	}
	return payload, err
}

func (codec *SessionCodec) Encode(payload []byte) (packet []byte, err error) {
	packet = make([]byte, 4+len(payload))

	// size
	binary.BigEndian.PutUint32(packet, uint32(len(payload)))

	if copy(packet[4:], payload) != len(payload) {
		err = errors.New("Encode: copy failed")
		return
	}
	return packet, err
}



const aesBlockLen = 16

type SessionAesCodec struct {
	r			io.Reader
	cipherBlk	cipher.Block
}

func NewSessionAesCodec(r io.Reader, key []byte) *SessionAesCodec {
	codec := &SessionAesCodec{
		r:			r,
	}
	cipherBlk, err := aes.NewCipher(key)
	if err != nil {
		glog.Errorf("NewSessionAesCodec: NewCipher(%d bytes) = %s", len(key), err)
		return nil
	}
	codec.cipherBlk = cipherBlk
	return codec
}

func (codec *SessionAesCodec) Decode() (payload []byte, err error) {
	firstBlock := make([]byte, aesBlockLen)
	_, err = io.ReadFull(codec.r, firstBlock)
	if err != nil {
		return nil, err
	}

	decFirstBlock := make([]byte, aesBlockLen)
	codec.cipherBlk.Decrypt(decFirstBlock, firstBlock)
	header := decFirstBlock[:4]
	length := binary.BigEndian.Uint32(header)
	if length > 64*1024*1024 {
		fmt.Errorf("Decode: invalid length %d", length)
		return
	}

	overrun := (4+length) % aesBlockLen
	paddingLen := aesBlockLen - overrun
	remainLen := 4+length+paddingLen-aesBlockLen
	if remainLen > 0 {
		remain := make([]byte, remainLen)
		_, err := io.ReadFull(codec.r, remain)
		if err != nil {
			return nil, err
		}
		decRemain := make([]byte, remainLen)
		for i := 0; i < int(remainLen/aesBlockLen); i++ {
			codec.cipherBlk.Decrypt(decRemain[i*aesBlockLen:(i+1)*aesBlockLen], remain[i*aesBlockLen:(i+1)*aesBlockLen])
		}

		payload = make([]byte, (aesBlockLen-4) + (remainLen-paddingLen))
		copy(payload, decFirstBlock[4:])
		copy(payload[(aesBlockLen-4):], decRemain[:remainLen-paddingLen])
		return payload, err 
	} else {
		payload = decFirstBlock[4:overrun]
		return payload, err
	}
}

func (codec *SessionAesCodec) Encode(payload []byte) (packet []byte, err error) {
	overrun := (4+len(payload)) % aesBlockLen
	paddingLen := aesBlockLen - overrun
	packetSize := 4+len(payload)+paddingLen
	buf := make([]byte, packetSize)

	binary.BigEndian.PutUint32(buf, uint32(len(payload)))
	copy(buf[4:], payload)

	packet = make([]byte, packetSize)
	for i := 0; i < packetSize/aesBlockLen; i++ {
		codec.cipherBlk.Encrypt(packet[i*aesBlockLen:(i+1)*aesBlockLen], buf[i*aesBlockLen:(i+1)*aesBlockLen])
	}
	return packet, nil 
}
