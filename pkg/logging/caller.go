package logging

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// 针对im-fun项目的Caller获取
func printCallerForImFun(f *runtime.Frame) string {
	pc, file, line, _ := runtime.Caller(9)
	dir, filename := path.Split(file)
	baseDir := path.Base(dir)
	// get function info
	sFun := runtime.FuncForPC(pc)
	s := strings.Split(sFun.Name(), ".")
	funcName := s[len(s)-1]

	return fmt.Sprintf("%v/%v:%v func[%v]", baseDir, filename, line, funcName)

	// 原来的机制：
	//s := strings.Split(f.Function, ".")
	//funcName := s[len(s)-1]
	//dir, filename := path.Split(f.File)
	//baseDir := path.Base(dir)
	//return fmt.Sprintf("%v/%v:%v:[%v]", baseDir, filename, f.Line, funcName)
}
