package client 

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
)

type Account struct {
	Name		string
	Passwd		string
	Zone		int32
}

type ClientConfig struct {
	CaFile				string
	LoginAddress		[]string
	Accounts			[]Account
}

func NewClientConfig(fileName string) *ClientConfig {
	clientConfig := &ClientConfig{}

	file, fileErr := ioutil.ReadFile(fileName)
	if fileErr != nil {
		glog.Errorf("load client configuration file %s failed %s", fileName, fileErr)
		return clientConfig
	}

	decoderErr := json.Unmarshal(file, clientConfig)
	if decoderErr != nil {
		glog.Errorf("decoding client configuration file %s failed %s", fileName, decoderErr)
	}

	return clientConfig
}

