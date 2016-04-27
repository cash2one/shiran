package masterslave

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
)

type MasterSlaveAppConfig struct {
	Name		string
	Bin			string
	Arg			[]string
}

type MasterSlaveConfig struct {
	Name		string
	App			[]MasterSlaveAppConfig
}

type MasterConfig struct {
	MasterAddress		string
	CommandAddress		string
	Slave				[]MasterSlaveConfig
}

func NewMasterConfig(fileName string) *MasterConfig {
	masterConfig := &MasterConfig{}

	file, fileErr := ioutil.ReadFile(fileName)
	if fileErr != nil {
		glog.Errorf("load master configuration file %s failed %s", fileName, fileErr)
		return masterConfig
	}

	decoderErr := json.Unmarshal(file, masterConfig)
	if decoderErr != nil {
		glog.Errorf("decoding master configuration file %s failed %s", fileName, decoderErr)
	}

	return masterConfig
}

