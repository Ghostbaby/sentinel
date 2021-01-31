package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/toolkits/file"
)

type GlobalConfig struct {
	IPS        map[string]IP `json:"ips"`
	Routines   int           `json:"routines"`
	RetryTimes int           `json:"retry_times"`
	IdRsa      string        `json:"id_rsa"`
	ScpFile    string        `json:"scp_file"`
}

type IP struct {
	Describe string `json:"describe"`
	Port     string `json:"port"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string, reload bool) {
	if cfg == "" {
		if reload {
			logrus.Error("configuration file is nil")
			return
		}
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		if reload {
			logrus.Error("config file:", cfg, "is not existent")
			return
		}
		log.Fatalln("config file:", cfg, "is not existent")
	}
	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		if reload {
			logrus.Error("read config file:", cfg, "fail:", err)
			return
		}
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		if reload {
			logrus.Error("parse config file:", cfg, "fail:", err)
			return
		}
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	if !reload {
		log.Println("read config file:", cfg, "successfully")
	}
}
