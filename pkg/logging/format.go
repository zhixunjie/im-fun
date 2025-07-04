package logging

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
)

// 自定义格式

type myFormatter struct {
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// time
	_, _ = fmt.Fprintf(b, "\x1b[%dm", colorDefault) // 添加颜色
	time := entry.Time.Format("2006-01-02 15:04:05.000")
	b.WriteString(time)

	// level
	levelColor := getColorByLevel(entry.Level)
	level := strings.ToUpper(entry.Level.String())
	_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor) // 添加颜色
	b.WriteString(" [" + level + "]")
	b.WriteString("\x1b[0m") // 去掉颜色

	// caller
	_, _ = fmt.Fprintf(b, "\x1b[%dm", colorYellow) // 添加颜色
	b.WriteString(" " + printCallerForImFun(entry.Caller))

	// msg
	_, _ = fmt.Fprintf(b, "\x1b[%dm", colorBlue) // 添加颜色
	b.WriteString(" ")
	b.WriteString(entry.Message)

	// other fields
	b.WriteString(" ")
	f.writeFields(b, entry)

	b.WriteByte('\n')
	b.WriteString("\x1b[0m") // 去掉颜色
	return b.Bytes(), nil
}

func (f *myFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *myFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
}

const (
	colorRed     = 31
	colorYellow  = 33
	colorBlue    = 36
	colorGray    = 37
	colorDefault = 97
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
