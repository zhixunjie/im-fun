package logging

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// printCaller 输出Caller信息
func printCaller(f *runtime.Frame) (string, string) {
	s := strings.Split(f.Function, ".")
	funcName := s[len(s)-1]
	dir, filename := path.Split(f.File)
	baseDir := path.Base(dir)
	return funcName, fmt.Sprintf("%v/%v:%v", baseDir, filename, f.Line)
}

// printCallerOther 输出Caller信息
func printCallerOther(f *runtime.Frame) string {
	s := strings.Split(f.Function, ".")
	funcName := s[len(s)-1]
	dir, filename := path.Split(f.File)
	baseDir := path.Base(dir)
	return fmt.Sprintf("%v/%v:%v:[%v]", baseDir, filename, f.Line, funcName)
}
