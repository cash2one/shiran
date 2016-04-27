package masterslave

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
)

type SlaveConfig struct {
	Name				string
	MasterAddress		string
}

func NewSlaveConfig(fileName string) *SlaveConfig {
	slaveConfig := &SlaveConfig{}

	file, fileErr := ioutil.ReadFile(fileName)
	if fileErr != nil {
		glog.Errorf("load slave configuration file %s failed %s", fileName, fileErr)
		return slaveConfig 
	}

	decoderErr := json.Unmarshal(file, slaveConfig)
	if decoderErr != nil {
		glog.Errorf("decoding slave configuration file %s failed %s", fileName, decoderErr)
	}

	return slaveConfig
}
