package masterslave

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
)

type CommanderConfig struct {
	MasterAddress		string
}

func NewCommanderConfig(fileName string) *CommanderConfig {
	commanderConfig := &CommanderConfig{}

	file, fileErr := ioutil.ReadFile(fileName)
	if fileErr != nil {
		glog.Errorf("load commander configuration file %s failed %s", fileName, fileErr)
		return commanderConfig 
	}

	decoderErr := json.Unmarshal(file, commanderConfig)
	if decoderErr != nil {
		glog.Errorf("decoding commander configuration file %s failed %s", fileName, decoderErr)
	}

	return commanderConfig 
}
