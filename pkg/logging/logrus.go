package logging

import (
	"github.com/sirupsen/logrus"
)

var loggerMyFormatter = logrus.New()

func init() {
	setLoggerMyFormatter(loggerMyFormatter)
}
func setLoggerMyFormatter(logger *logrus.Logger) {
	logger.SetReportCaller(true)
	logger.SetFormatter(&myFormatter{})
}
