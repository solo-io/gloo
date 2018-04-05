package log

import (
	"fmt"
	"os"
	"runtime"

	"regexp"

	"io"

	"time"

	"github.com/k0kubun/pp"
)

var debugMode = os.Getenv("DEBUG") == "1"
var DefaultOut io.Writer = os.Stdout

func init() {
	if os.Getenv("DISABLE_COLOR") == "1" {
		pp.ColoringEnabled = false
	}
}

var rxp = regexp.MustCompile(".*/src/")

func Sprintf(format string, a ...interface{}) string {
	return pp.Sprintf("%v\t"+format, append([]interface{}{line()}, a...)...)
}

func GreyPrintf(format string, a ...interface{}) {
	fmt.Fprintf(DefaultOut, "%v\t"+format+"\n", append([]interface{}{line()}, a...)...)
}

func Printf(format string, a ...interface{}) {
	pp.Fprintf(DefaultOut, "%v\t"+format+"\n", append([]interface{}{line()}, a...)...)
}

func Warnf(format string, a ...interface{}) {
	pp.Fprintf(DefaultOut, "WARNING: %v\t"+format+"\n", append([]interface{}{line()}, a...)...)
}

func Debugf(format string, a ...interface{}) {
	if debugMode {
		pp.Fprintf(DefaultOut, "%v\t"+format+"\n", append([]interface{}{line()}, a...)...)
	}
}

func Fatalf(format string, a ...interface{}) {
	pp.Fprintf(DefaultOut, "%v\t"+format+"\n", append([]interface{}{line()}, a...)...)
	os.Exit(1)
}

func line() string {
	_, file, line, _ := runtime.Caller(2)
	file = rxp.ReplaceAllString(file, "")
	return fmt.Sprintf("%v: %v:%v", time.Now().Format(time.RFC1123), file, line)
}
