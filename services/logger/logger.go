package logger

import (
	"app/config"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

var once sync.Once

func Initialize() {
	once.Do(func() {
		_conf := config.GetConfig().Others
		if _conf.LogFormat == "json" {
			log.SetFormatter(&log.JSONFormatter{
				TimestampFormat: "02-01-2006 15:04:05",
			})
		} else {
			log.SetFormatter(&log.TextFormatter{
				ForceColors:     true,
				FullTimestamp:   true,
				TimestampFormat: "02-01-2006 15:04:05",
			})
		}
		if level, err := log.ParseLevel(_conf.LogLevel); err == nil {
			log.SetLevel(level)
			log.Infoln("log level set to ", _conf.LogLevel)
		}

		log.SetReportCaller(true)
		log.SetOutput(os.Stdout)
	})
}
