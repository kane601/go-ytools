package tools

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// GetSysBit 获取当前位数
func GetSysBit() int {
	if "amd64" == runtime.GOARCH {
		return 64
	}
	return 32
}

// GetSpecialDir 获取指定目录
func GetSpecialDir(t SpecialDirType) string {
	home := GetEnv("HOME")
	switch t {
	case LocalAppdata:
		return AbsJoinPath(home, "Library/Application Support")
	case LocalRoamingAppdata:
		return AbsJoinPath(home, "Library/Application Support")
	case TempLocation:
		return AbsJoinPath(home, "Library/Caches")
	case PreferencesLocation:
		return AbsJoinPath(home, "Library/Preferences")
	case HomeLocation:
		return home
	default:
	}
	return home
}

// SendCtrlC
func SendCtrlC(p *os.Process) bool {
	return nil == p.Signal(os.Interrupt)
}
func SendCtrlC2(p *os.Process) bool {
	return nil == p.Signal(os.Interrupt)
}

func SendCtrlCWait(p *os.Process) bool {
	return SendCtrlC(p)
}
func SendCtrlCWait2(p *os.Process) bool {
	return SendCtrlC(p)
}

func SendCtrlBreak(p *os.Process) bool {
	return nil == p.Signal(os.Interrupt)
}
func SendCtrlBreak2(p *os.Process) bool {
	return nil == p.Signal(os.Interrupt)
}

func SendCtrlC64(p *os.Process) bool {
	return SendCtrlC(p)
}
func SendCtrlC642(p *os.Process) bool {
	return SendCtrlC(p)
}

func SendCtrlCWait64(p *os.Process) bool {
	return SendCtrlCWait(p)
}
func SendCtrlCWait642(p *os.Process) bool {
	return SendCtrlCWait(p)
}

func SendCtrlBreak64(p *os.Process) bool {
	return SendCtrlBreak(p)
}
func SendCtrlBreak642(p *os.Process) bool {
	return SendCtrlBreak(p)
}

func CommandNoWin(cmd string, args ...string) *exec.Cmd {
	return exec.Command(cmd, args...)
}

func GetAnotherInstancePID() int {
	out, err := exec.Command("ps", "-eo", "pid,args").Output()
	if err != nil {
		return 0
	}

	processArgs := strings.Split(string(out), "\n")
	if len(processArgs) == 0 {
		return 0
	}
	for _, processArg := range processArgs[1:] {
		processInfo := strings.Fields(processArg)
		if len(processInfo) == 0 {
			continue
		}
		if !strings.Contains(processArg, os.Args[0]) {
			continue
		}
		pid, _ := strconv.Atoi(processInfo[0])
		if os.Getpid() == pid {
			continue
		}
		return pid
	}
	return 0
}

func GetPIDByParent(path string, ppid int) int {
	return 0
}

func WaitPID(pid int, _ time.Duration) {
	pro, err := os.FindProcess(int(pid))
	if err == nil {
		pro.Wait()
		return
	}
}

func IsWin7() bool {
	return false
}
