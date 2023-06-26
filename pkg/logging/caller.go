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

func printCallerForImFun(f *runtime.Frame) string {
	pc, file, line, _ := runtime.Caller(9)
	dir, filename := path.Split(file)
	baseDir := path.Base(dir)
	// get function info
	sFun := runtime.FuncForPC(pc)
	s := strings.Split(sFun.Name(), ".")
	funcName := s[len(s)-1]

	return fmt.Sprintf("%v/%v:%v:[%v]", baseDir, filename, line, funcName)

	// 原来的机制：
	//s := strings.Split(f.Function, ".")
	//funcName := s[len(s)-1]
	//dir, filename := path.Split(f.File)
	//baseDir := path.Base(dir)
	//return fmt.Sprintf("%v/%v:%v:[%v]", baseDir, filename, f.Line, funcName)
}
