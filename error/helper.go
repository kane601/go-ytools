package yerror

import (
	"regexp"
	"runtime"
	"strings"

	tools "github.com/kane601/go-ytools/utils"
)

// CallerInfo ...
type CallerInfo struct {
	File string `json:"file"`
	Fun  string `json:"name"`
	Line int    `json:"line"`
}

var (
	reInit    = regexp.MustCompile(`init·?\d+$`)
	reClosure = regexp.MustCompile(`func·?\d+$`)
)

// Caller 获取调用信息
func Caller(skip int) CallerInfo {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return CallerInfo{
			File: "???",
			Fun:  "???",
			Line: -1,
		}
	}
	fun := runtime.FuncForPC(pc).Name()
	if reInit.MatchString(fun) {
		fun = reInit.ReplaceAllString(fun, "init")
	} else if reClosure.MatchString(fun) {
		fun = reClosure.ReplaceAllString(fun, "func")
	}
	file = tools.PathName(file)
	return CallerInfo{
		File: file,
		Fun:  fun,
		Line: line,
	}
}

// CallerList 获取调用信息链
func CallerList(skip int) (ret []CallerInfo) {
	for ; ; skip++ {
		info := Caller(skip + 1)
		if info.Line == -1 || strings.HasPrefix(info.Fun, "runtime.") {
			break
		}
		ret = append(ret, info)
	}
	return
}
