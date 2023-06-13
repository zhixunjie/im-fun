package logging

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// PrintCaller 输出文件
func PrintCaller(f *runtime.Frame) (string, string) {
	s := strings.Split(f.Function, ".")
	funcName := s[len(s)-1]
	dir, filename := path.Split(f.File)
	baseDir := path.Base(dir)
	return funcName, fmt.Sprintf("%v/%v:%v", baseDir, filename, f.Line)
}

// PrintCallerOther 输出文件
func PrintCallerOther(f *runtime.Frame) string {
	s := strings.Split(f.Function, ".")
	funcName := s[len(s)-1]
	dir, filename := path.Split(f.File)
	baseDir := path.Base(dir)
	return fmt.Sprintf("%v/%v:%v:[%v]", baseDir, filename, f.Line, funcName)
}
