package tools

import (
	"os"
	"strings"
)

// GetEnv 获取环境变量
func GetEnv(key string) string {
	for _, item := range os.Environ() {
		keyvals := strings.Split(item, "=")
		if keyvals[0] == key {
			return keyvals[1]
		}
	}
	return ""
}

func OneImageInstance() (isKill bool) {
	pid := GetAnotherInstancePID()
	if pid != 0 {
		isKill = true
		SendCtrlCWait(&os.Process{
			Pid: pid,
		})
	}
	return
}

func GetSubProcessPID(name string, ppid int) (subpid int) {
	if name == "" || ppid <= 0 {
		return
	}
	subpid = GetPIDByParent(name, ppid)
	return
}

func KillSubProcess(name string, ppid int, is64 bool) {
	pid := GetPIDByParent(name, ppid)
	if pid > 0 {
		if is64 {
			SendCtrlC(&os.Process{Pid: pid})
		} else {
			SendCtrlC64(&os.Process{Pid: pid})
		}
	}
}
