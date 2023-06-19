package yaml

import (
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func LoadConfig(path string, conf interface{}) (err error) {
	logHead := fmt.Sprintf("LoadConfig path=%v", path)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logging.Errorf("err=%v", err)
		return err
	}

	// begin to unmarshal
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return err
	}
	logging.Infof(logHead+"config=%+v", conf)

	return nil
}
