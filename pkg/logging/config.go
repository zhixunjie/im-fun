package logging

import (
	"github.com/sirupsen/logrus"
)

var LoggingMyFormatter = logrus.New()

func InitLogConfig() {
	SetLoggingMyFormatter(LoggingMyFormatter)
}

func SetLoggingMyFormatter(logger *logrus.Logger) {
	logger.SetReportCaller(true)
	logger.SetFormatter(&MyFormatter{})
}

func Info(args ...interface{}) {
	LoggingMyFormatter.Info(args...)
}

func Infof(format string, args ...interface{}) {
	LoggingMyFormatter.Infof(format, args...)
}

func Debug(args ...interface{}) {
	LoggingMyFormatter.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	LoggingMyFormatter.Debugf(format, args...)
}

func Error(args ...interface{}) {
	LoggingMyFormatter.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	LoggingMyFormatter.Errorf(format, args...)
}

func Warn(args ...interface{}) {
	LoggingMyFormatter.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	LoggingMyFormatter.Warnf(format, args...)
}
