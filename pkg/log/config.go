package log

import (
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
	"time"
)

func InitLogConfig() {
	logrus.SetFormatter(&logrus.TextFormatter{
		//DisableColors: true,
		FullTimestamp: true,
	})
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcName, filename
		},
	})
}
