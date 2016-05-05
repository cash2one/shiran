#!/bin/bash

#./easyrsa init-pki
#./easyrsa build-ca
#./easyrsa --subject-alt-name=DNS:*.a.com,DNS:*.b.com,DNS:*.c.com build-server-full *.a.com nopass

USERPROTO_PROTO_PATH=~/goyard/src/github.com/williammuji/shiran/proto/userproto/
USERPROTO_GO_OUT=~/goyard/src/github.com/williammuji/shiran/proto/userproto/

GATELOGIN_PROTO_PATH=~/goyard/src/github.com/williammuji/shiran/proto/gatelogin/
GATELOGIN_GO_OUT=~/goyard/src/github.com/williammuji/shiran/proto/gatelogin/

MASTERSLAVE_PROTO_PATH=~/goyard/src/github.com/williammuji/shiran/masterslave
MASTERSLAVE_GO_OUT=~/goyard/src/github.com/williammuji/shiran/masterslave

protoc --proto_path=${USERPROTO_PROTO_PATH} --go_out=${USERPROTO_GO_OUT} ${USERPROTO_PROTO_PATH}/*.proto
protoc --proto_path=${GATELOGIN_PROTO_PATH} --go_out=${GATELOGIN_GO_OUT} ${GATELOGIN_PROTO_PATH}/*.proto
protoc --proto_path=${MASTERSLAVE_PROTO_PATH} --go_out=${MASTERSLAVE_GO_OUT} ${MASTERSLAVE_PROTO_PATH}/*.proto

go install github.com/williammuji/shiran/proto/userproto
go install github.com/williammuji/shiran/proto/gatelogin
go install github.com/williammuji/shiran/shiran
go install github.com/williammuji/shiran/examples/pingpong/pingpongserver/pingpongserver
go install github.com/williammuji/shiran/examples/pingpong/pingpongclient/pingpongclient
go install github.com/williammuji/shiran/masterslave
go install github.com/williammuji/shiran/master
go install github.com/williammuji/shiran/slave
go install github.com/williammuji/shiran/commander
go install github.com/williammuji/shiran/login
go install github.com/williammuji/shiran/login/login
go install github.com/williammuji/shiran/gate
go install github.com/williammuji/shiran/gate/gate
go install github.com/williammuji/shiran/client
go install github.com/williammuji/shiran/client/client
