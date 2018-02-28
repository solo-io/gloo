package log

import (
	"fmt"
	"os"
	"runtime"

	"github.com/k0kubun/pp"
)

var debugMode = os.Getenv("DEBUG") == "1"

func Sprintf(format string, a ...interface{}) string {
	_, file, line, _ := runtime.Caller(1)
	return pp.Sprintf("%v:%v\n"+format, append([]interface{}{file, line}, a...)...)
}

func GreyPrintf(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("%v:%v\n"+format+"\n", append([]interface{}{file, line}, a...)...)
}

func Printf(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	pp.Printf("%v:%v\n"+format+"\n", append([]interface{}{file, line}, a...)...)
}

func Warnf(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	pp.Printf("WARNING: %v:%v\n"+format+"\n", append([]interface{}{file, line}, a...)...)
}

func Debugf(format string, a ...interface{}) {
	if debugMode {
		_, file, line, _ := runtime.Caller(1)
		pp.Printf("%v:%v\n"+format+"\n", append([]interface{}{file, line}, a...)...)
	}
}

func Fatalf(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	pp.Fatalf("%v:%v\n"+format+"\n", append([]interface{}{file, line}, a...)...)
}
