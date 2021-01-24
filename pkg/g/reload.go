package g

import (
	"time"

	"github.com/go-logr/logr"
)

func reloadConfig(log logr.Logger) {
	ParseConfig(ConfigFile, true)
	log.V(4).Info("reload config complete")

}

func ConfigReload(log logr.Logger) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		reloadConfig(log)
		<-ticker.C
	}
}
