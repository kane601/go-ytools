package tools

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/kane601/go-w32"
	"github.com/kane601/go-w32/wutil"
)

// GetSysBit 获取当前位数
func GetSysBit() int {
	return w32.GetSysBit()
}

// GetSpecialDir 获取指定目 录
func GetSpecialDir(t SpecialDirType) string {
	switch t {
	case TempLocation:
		return os.TempDir()
	case PreferencesLocation:
		return AbsJoinPath(GetSpecialDir(LocalRoamingAppdata), "Preferences")
	case HomeLocation:
		return AbsJoinPath(w32.SHGetSpecialFolderPath(int32(0x0010)), "..")
	default:
	}
	return w32.SHGetSpecialFolderPath(int32(t))
}

func CommandNoWin(cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	return c
}

func SendCtrlC(p *os.Process) bool {
	return sendCtrlC(p, false, true)
}

func SendCtrlC2(p *os.Process) bool {
	return sendCtrlC(p, false, false)
}

func SendCtrlC64(p *os.Process) bool {
	return sendCtrlC(p, true, true)
}

func SendCtrlC642(p *os.Process) bool {
	return sendCtrlC(p, true, false)
}

func SendCtrlCWait(p *os.Process) bool {
	return sendCtrlCWait(p, false, true)
}
func SendCtrlCWai2(p *os.Process) bool {
	return sendCtrlCWait(p, false, false)
}

func SendCtrlCWait64(p *os.Process) bool {
	return sendCtrlCWait(p, true, true)
}
func SendCtrlCWait642(p *os.Process) bool {
	return sendCtrlCWait(p, true, false)
}

func SendCtrlBreak(p *os.Process) bool {
	return SendCtrlC(p)
}

func SendCtrlBreak2(p *os.Process) bool {
	return SendCtrlC2(p)
}

func SendCtrlBreak64(p *os.Process) bool {
	return SendCtrlC64(p)
}

func SendCtrlBreak642(p *os.Process) bool {
	return SendCtrlC642(p)
}

func GetAnotherInstancePID() int {
	pid := 0
	wutil.ProcessWalk(func(pro w32.PROCESSENTRY32) bool {
		if wutil.IsMatchProcess(os.Args[0], pro) && pro.Th32ProcessID != w32.DWORD(os.Getpid()) {
			pid = int(pro.Th32ProcessID)
			return false
		}
		return true
	})
	return pid
}

func sendCtrlC(p *os.Process, isWindowsKill64 bool, nowin bool) bool {
	name := "./windows-kill"
	if isWindowsKill64 {
		name = "./windows-kill64"
	}
	var c *exec.Cmd
	if nowin {
		c = CommandNoWin(name, "-SIGBREAK", fmt.Sprintf("%d", p.Pid))
	} else {
		c = exec.Command(name, "-SIGBREAK", fmt.Sprintf("%d", p.Pid))
	}
	if e := c.Start(); e != nil {
		fmt.Println("sendCtrlC:", e)
		return false
	}
	return true
}

func sendCtrlCWait(p *os.Process, isWindowsKill64 bool, nowin bool) bool {
	name := "windows-kill"
	if isWindowsKill64 {
		name = "windows-kill64"
	}
	var c *exec.Cmd
	if nowin {
		c = CommandNoWin(name, "-SIGBREAK", fmt.Sprintf("%d", p.Pid))
	} else {
		c = exec.Command(name, "-SIGBREAK", fmt.Sprintf("%d", p.Pid))
	}

	if e := c.Start(); e != nil {
		fmt.Println("sendCtrlCWait:", e)
		return false
	}
	c.Wait()
	return true
}

func GetPIDByParent(path string, ppid int) int {
	return wutil.GetPIDByParent(path, ppid)
}

func WaitPID(pid int, inter time.Duration) {
	pro, err := os.FindProcess(int(pid))
	if err == nil {
		pro.Wait()
		return
	}

	for {
		if !isProcessRuningPid(int64(pid)) {
			return
		}
		time.Sleep(inter)
	}
}

func isProcessRuningPid(pid int64) (bret bool) {
	if pid <= 0 {
		return
	}
	wutil.ProcessWalk(func(e32 w32.PROCESSENTRY32) bool {
		if int64(e32.Th32ProcessID) == int64(pid) {
			bret = true
			return false
		}
		return true
	})
	return
}

var _isWin7 int = -1

func IsWin7() bool {
	if _isWin7 != -1 {
		return _isWin7 == 1
	}
	m, j, _ := w32.GetSystemVersion()
	if m == 6 && j == 1 {
		_isWin7 = 1
	} else {
		_isWin7 = 0
	}
	return _isWin7 == 1
}
