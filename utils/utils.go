package tools

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

func URLProtol(s string) string {
	pu, err := url.Parse(s)
	if err == nil && pu.Scheme != "" {
		return pu.Scheme
	}
	return ""
}

func GetSystemProxy() (proxy string) {
	defer func() {
		if proxy == "" {
			return
		}
		pu, err := url.Parse(proxy)
		if err == nil && pu.Scheme != "" {
			return
		}
		proxy = "http://" + proxy
	}()
	if strings.EqualFold(runtime.GOOS, "windows") {
		enable := false
		cmd := exec.Command("reg", "query", `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyEnable")
		by, err := cmd.Output()
		if err != nil {
			return
		}
		s := strings.TrimSpace(string(by))
		reEnable := regexp.MustCompile(`ProxyEnable\s+REG_DWORD\s+(.+)`)
		if reEnable.Match(by) {
			arr := reEnable.FindStringSubmatch(s)
			if len(arr) == 2 {
				enable = arr[1] == "1" || arr[1] == "0x1"
			}
		} else {
			return
		}
		if !enable {
			return
		}

		cmd = exec.Command("reg", "query", `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyServer")
		by, err = cmd.Output()
		if err != nil {
			return
		}
		s = strings.TrimSpace(string(by))
		reProxy := regexp.MustCompile(`ProxyServer\s+REG_SZ\s+(.+)`)
		if reProxy.Match(by) {
			arr := reProxy.FindStringSubmatch(s)
			if len(arr) == 2 {
				proxy = strings.TrimSpace(arr[1])
			}
		} else {
			return
		}

		if strings.Contains(s, "=") {
			for _, a := range strings.Split(s, " ") {
				if strings.Contains(a, "=") {
					s = a
					break
				}
			}

			for _, p := range strings.Split(s, ";") {
				addr := strings.Split(p, "=")
				if len(addr) == 1 && URLProtol(addr[0]) != "" {
					proxy = addr[0]
					break
				}
				if len(addr) == 2 && URLProtol(addr[1]) == "" && (strings.HasPrefix(addr[0], "http") || strings.HasPrefix(addr[0], "socks")) {
					proxy = addr[0] + "://" + addr[1]
					break
				}
				if len(addr) == 2 && URLProtol(addr[1]) != "" {
					proxy = addr[1]
					break
				}
			}
		}
	} else {
		cmd := exec.Command("scutil", "--proxy")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return
		}
		cmd.Start()
		r := bufio.NewReader(stdout)
		kv := map[string]string{}
		for {
			l, _, e := r.ReadLine()
			if e != nil {
				break
			}
			line := strings.TrimSpace(string(l))
			if !strings.Contains(line, ":") {
				continue
			}
			lineArr := strings.Split(line, ":")
			if len(lineArr) != 2 {
				continue
			}
			kv[strings.TrimSpace(lineArr[0])] = strings.TrimSpace(lineArr[1])
		}
		cmd.Wait()
		if enable, ok := kv["SOCKSEnable"]; ok && enable == "1" {
			proxy = "socks5://" + kv["SOCKSProxy"] + ":" + kv["SOCKSPort"]
			return
		}
		if enable, ok := kv["HTTPSEnable"]; ok && enable == "1" {
			proxy = "https://" + kv["HTTPSProxy"] + ":" + kv["HTTPSPort"]
			return
		}
		if enable, ok := kv["HTTPEnable"]; ok && enable == "1" {
			proxy = "http://" + kv["HTTPProxy"] + ":" + kv["HTTPPort"]
			return
		}

	}
	return
}

//VersionSplit 拆分版本信息
func VersionSplit(str string) []string {
	reg := regexp.MustCompile(`\d+`)
	vers := reg.FindAllString(str, -1)
	var results []string
	for _, ver := range vers {
		ver = strings.TrimLeft(ver, "0")
		if ver == "" {
			ver = "0"
		}
		if len(ver) > 0 {
			results = append(results, ver)
		}
	}
	return results
}

//CmpVersion 比较版本
func CmpVersion(str1 string, str2 string) int {
	results1 := VersionSplit(str1)
	results2 := VersionSplit(str2)

	if len(results1) > len(results2) {
		for i := len(results1) - len(results2); i > 0; i-- {
			results2 = append(results2, "0")
		}
	} else if len(results1) < len(results2) {
		for i := len(results2) - len(results1); i > 0; i-- {
			results1 = append(results1, "0")
		}
	}
	for i := 0; i < len(results1); i++ {
		if len(results1[i]) != len(results2[i]) {
			if len(results1[i]) > len(results2[i]) {
				return 1
			}
			return -1
		} else if 0 != strings.Compare(results1[i], results2[i]) {
			return strings.Compare(results1[i], results2[i])
		}
	}
	return 0
}

//IsInArray ...
func IsInArray(arr, val interface{}) bool {
	arrValueOf := reflect.ValueOf(arr)
	for i := 0; i < arrValueOf.Len(); i++ {
		if reflect.DeepEqual(arrValueOf.Index(i).Interface(), val) {
			return true
		}
	}
	return false
}

//ReplaceString 根据map替换文本内容到
func ReplaceString(strdata string, m map[string]string) string {
	for key, val := range m {
		strdata = strings.ReplaceAll(strdata, key, val)
	}
	return strdata
}

//AttachWait ...
func AttachWait() {
	fmt.Println(os.Getpid())
	<-time.After(time.Second * 30)
}

//CmpDate ...
func CmpDate(t1, t2 time.Time) int64 {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	tt1 := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	tt2 := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	return tt1.Unix() - tt2.Unix()
}

func testValue(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expect:%v,got:%v", want, got)
	}
}

func _newStopInputCtx(dur time.Duration, binput bool) context.Context {
	var ctx context.Context
	var cancleFun context.CancelFunc
	if dur <= 0 {
		ctx, cancleFun = context.WithCancel(context.Background())
	} else {
		ctx, cancleFun = context.WithTimeout(context.Background(), dur)
	}

	if binput {
		go func() {
			r := bufio.NewReader(os.Stdin)
			line, _, _ := r.ReadLine()
			sline := strings.TrimSpace(string(line))
			fmt.Println("recv input:", sline)
			if strings.EqualFold(sline, "stop") {
				fmt.Println("cancleFun()")
				cancleFun()
			}
		}()
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-ch
		cancleFun()
	}()
	return ctx
}

func NewStopCtx() context.Context {
	return _newStopInputCtx(0, false)
}

func NewStopCtxWithTimeout(dur time.Duration) context.Context {
	return _newStopInputCtx(dur, false)
}

func NewStopInputCtx() context.Context {
	return _newStopInputCtx(0, true)
}

func NewStopInputCtxWithTimeout(dur time.Duration) context.Context {
	return _newStopInputCtx(dur, true)
}

func SelectFile(name string, ignoreLocal bool) string {
	path := TempPath(name)
	if IsExist(path) && !IsExist(path+".lock") {
		return path
	}
	path = DataPath(name)
	if IsExist(path) && !IsExist(path+".lock") {
		return path
	}
	if !ignoreLocal {
		path = LocalPath(name)
		if IsExist(path) {
			return path
		}
	}
	return ""
}

func AutoSelectAndCopyFile(name string, x, notLocal bool) string {
	path := SelectFile(name, notLocal)
	if len(path) == 0 {
		path = TempPath(name)
		var err error
		if IsExist(LocalPath(name)) {
			err = CopyFile(context.TODO(), LocalPath(name), path)
		} else {
			err = CopyFile(context.TODO(), LocalPath(name+"_bk"), path)
		}
		if nil == err {
			if !strings.EqualFold(runtime.GOOS, "windows") && x {
				err = exec.Command("chmod", "a+x", path).Run()
				fmt.Println("chmod a+x:", err)
			}
		} else {
			path = ""
			fmt.Println("AutoSelectAndCopyFile: ", err)
		}
	}
	if path == "" {
		path = LocalPath(name)
	}
	return path
}

func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		break
	}
	return false
}

func GetProcessExitCode(err error) (code int, ok bool) {
	if err == nil {
		return 0, true
	}
	exiterr, ok := err.(*exec.ExitError)
	if ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			code = status.ExitStatus()
		}
	}
	return
}

func CallWithProgress(ctx context.Context, doF, cancleF func(), progPerSecond float64, progF func(prog float64)) {
	if progPerSecond == 0 {
		progPerSecond = 1
	}

	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	ch := make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		prog := float64(0)
	loop:
		for {
			select {
			case <-ch:
				break loop
			case <-ctx.Done():
				if cancleF != nil {
					cancleF()
				}
				break loop
			default:
				prog += progPerSecond / 2
				if prog >= 95 {
					prog = 95
				}
				if progF != nil {
					progF(prog)
				}
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	if doF != nil {
		doF()
	}
	if !IsContextDone(ctx) && progF != nil {
		progF(100)
	}
	ch <- struct{}{}
}

func mapArrPathVal(m interface{}, paths []string) (ret interface{}, err error) {

	defer func() {
		recover()
	}()

	if len(paths) == 0 {
		ret = m
		return
	}
	vo := reflect.ValueOf(m)
	k := paths[0]
	kind := vo.Kind()
	if reflect.Array == kind || reflect.Slice == kind {
		i, e := strconv.ParseInt(k, 0, 63)
		if e != nil {
			err = e
			return
		}
		if vo.Len() <= int(i) || vo.Len() == 0 {
			err = fmt.Errorf("index too long")
			return
		}
		if i == -1 {
			i = int64(vo.Len() - 1)
		}
		return mapArrPathVal(vo.Index(int(i)).Interface(), paths[1:])
	} else if reflect.Map == vo.Kind() {
		return mapArrPathVal(vo.MapIndex(reflect.ValueOf(k)).Interface(), paths[1:])
	} else {
		err = fmt.Errorf("type is error")
	}
	return
}

func MapArrPathVal(m interface{}, path string) (ret interface{}, err error) {
	return mapArrPathVal(m, strings.Split(path, "/"))
}

func FindJSONString(data string, obj bool) (ret string) {
	startKey := '{'
	endKey := '}'
	if !obj {
		startKey = '['
		endKey = ']'
	}

	s := strings.Index(data, string(startKey))
	if s == -1 {
		return
	}
	data = data[s:]
	num := 0
	for i, c := range data {
		if c == rune(startKey) {
			num++
		} else if c == rune(endKey) {
			num--
		}
		if num == 0 {
			ret = data[:i+1]
			break
		}
	}
	return
}
