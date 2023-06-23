package logging

// 方便调用的函数

func Info(args ...interface{}) {
	loggerMyFormatter.Info(args...)
}

func Infof(format string, args ...interface{}) {
	loggerMyFormatter.Infof(format, args...)
}

func Debug(args ...interface{}) {
	loggerMyFormatter.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	loggerMyFormatter.Debugf(format, args...)
}

func Error(args ...interface{}) {
	loggerMyFormatter.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	loggerMyFormatter.Errorf(format, args...)
}

func Warn(args ...interface{}) {
	loggerMyFormatter.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	loggerMyFormatter.Warnf(format, args...)
}
