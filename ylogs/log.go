package ylogs

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sync"

	js "github.com/bitly/go-simplejson"
	yerror "github.com/kane601/go-ytools/error"
	tools "github.com/kane601/go-ytools/utils"
)

type Fmter interface {
	Println(a ...any) error
	Printf(format string, a ...any) error
}

// Loger loger接口
type Loger interface {

	//DbgPrintf ...
	DbgPrintf(format string, v ...interface{})

	//DbgPrint ...
	DbgPrint(v ...interface{})

	//TracePrintf ...
	TracePrintf(format string, v ...interface{})

	//TracePrint ...
	TracePrint(v ...interface{})

	//StdPrintf ...
	StdPrintf(format string, v ...interface{})

	//StdPrint ...
	StdPrint(v ...interface{})

	//CodePrint ...
	CodePrint(e interface{}, v ...interface{})

	//ProgressPrint 标准输出进度相关信息
	ProgressPrint(phase, progress, speed, size string)

	//ProgressPrintWithValue 标准输出进度相关信息
	ProgressPrintWithValue(val map[string]interface{}, progress string)
}

// New ...
func New() Loger {
	return newLoger("", nil, false)
}

func NewWithLocker() Loger {
	return newLoger("", nil, true)
}

func NewCodeTrace() Loger {
	l := newLoger("", nil, false)
	l.codeTrace = true
	return l
}

func NewWithFmter(fmter Fmter, lock ...bool) Loger {
	if len(lock) > 0 {
		return newLoger("", fmter, lock[0])
	}
	return newLoger("", fmter, false)
}

// NewWithFile ...
func NewWithFile(file string) Loger {
	return newLoger(file, nil, true)
}

// NewNull ...
func NewNull() Loger {
	return &NullLog{}
}

// NullLog 空日志
type NullLog struct {
}

// DbgPrintf ...
func (p *NullLog) DbgPrintf(format string, v ...interface{}) {}

// DbgPrint ...
func (p *NullLog) DbgPrint(v ...interface{}) {}

// TracePrintf ...
func (p *NullLog) TracePrintf(format string, v ...interface{}) {}

// TracePrint ...
func (p *NullLog) TracePrint(v ...interface{}) {}

// StdPrintf ...
func (p *NullLog) StdPrintf(format string, v ...interface{}) {}

// StdPrint ...
func (p *NullLog) StdPrint(v ...interface{}) {}

// CodePrint ...
func (p *NullLog) CodePrint(e interface{}, v ...interface{}) {}

func (p *NullLog) ProgressPrint(phase, progress, speed, size string) {}

func (p *NullLog) ProgressPrintWithValue(val map[string]interface{}, progress string) {}

// Log ...
type Log struct {
	log.Logger
	fmter     Fmter
	codeTrace bool
	mu        *sync.Mutex
}

func newLoger(file string, fmter Fmter, lock bool) *Log {
	l := new(Log)
	if fmter == nil {
		fmter = &fmterDefault{}
	}
	if lock {
		l.mu = &sync.Mutex{}
	}
	l.fmter = fmter
	var w io.Writer
	if len(file) > 0 {
		f, err := tools.OpenApptendFile(file)
		if err != nil {
			l._fmtPrintf("Log create fail:%v\n", err)
			w = os.Stdout
		} else {
			w = io.MultiWriter(os.Stdout, f)
		}
	} else {
		w = os.Stdout
	}
	l.SetOutput(w)
	l.SetFlags(log.LstdFlags)
	return l
}

// DbgPrintf ...
func (p *Log) DbgPrintf(format string, v ...interface{}) {
	p.print(odbg, format, v...)
}

// DbgPrint ...
func (p *Log) DbgPrint(v ...interface{}) {
	p.print(odbg, "", v...)
}

// TracePrintf ...
func (p *Log) TracePrintf(format string, v ...interface{}) {
	p.print(otrace, format, v...)
}

// TracePrint ...
func (p *Log) TracePrint(v ...interface{}) {
	p.print(otrace, "", v...)
}

// StdPrintf ...
func (p *Log) StdPrintf(format string, v ...interface{}) {
	p.print(ostd, format, v...)
}

// StdPrint ...
func (p *Log) StdPrint(v ...interface{}) {
	p.print(ostd, "", v...)
}

// CodePrint ...
func (p *Log) CodePrint(e interface{}, v ...interface{}) {
	var (
		code     int
		msg      string
		data     interface{}
		callinfo string
	)
	if e == nil {
		code = 0
	} else if c, ok := e.(int); ok {
		code = c
		msg = yerror.GetCodeTranslate(c)
	} else if yerr, ok := e.(yerror.Error); ok {
		code = yerr.Code()
		msg = yerr.Error()
		callinfo = yerr.CallerInfoStr()
	} else if err, ok := e.(error); ok {
		code = -1
		msg = err.Error()
	} else {
		code = -1
		msg = fmt.Sprint(e)
	}

	vt := make([]interface{}, 0, len(v))
	for _, i := range v {
		if i == nil {
			continue
		}
		tv := reflect.TypeOf(i)
		k := tv.Kind()
		if k == reflect.Array || k == reflect.Slice || k == reflect.Map || k == reflect.Struct {
			data = i
		} else {
			vt = append(vt, i)
		}
	}
	if len(vt) > 0 {
		if len(msg) > 0 {
			msg += " | "
		}
		msg += fmt.Sprint(vt...)
	}
	p.codePrint(code, msg, data, p.codeTrace, callinfo)

}

// ProgressPrint 标准输出进度相关信息
func (p *Log) ProgressPrintWithValue(val map[string]interface{}, progress string) {
	data := js.New()
	for k, v := range val {
		data.Set(k, v)
	}
	prog, ok := val["progress"]
	if !ok {
		data.Set("progress", correctProg100(progress))
	} else if p, ok := prog.(string); ok {
		data.Set("progress", correctProg100(p))
	}
	p.codePrint(1, "", data, false, "")
}

func (p *Log) ProgressPrint(phase, progress, speed, size string) {
	data := js.New()

	data.Set("phase", phase)
	data.Set("progress", correctProg100(progress))
	data.Set("speed", speed)
	data.Set("size", size)

	p.codePrint(1, "", data, false, "")
}

func (p *Log) codePrint(code int, msg string, data interface{}, record bool, stack string) {
	//for ui
	p._fmtPrintln(createJSONLog(code, msg, data))

	//for other
	if record {
		if len(stack) > 0 {
			if len(msg) > 0 {
				p.print(otrace, "%s & code(%d),stack(%s)", msg, code, stack)
			} else {
				p.print(otrace, "code(%d),stack(%s)", code, stack)
			}
		} else {
			if len(msg) > 0 {
				p.print(otrace, "%s & code(%d)", msg, code)
			} else {
				p.print(otrace, "code(%d)", code)
			}

		}

	}
}

// print ...
func (p *Log) print(lev int, format string, v ...interface{}) {
	p.SetPrefix(otoPrefix(lev))
	s := ""
	if ostd == lev {
		p._fmtPrintln(v...)
		return
	} else if len(format) > 0 {
		s = fmt.Sprintf(format, v...)
	} else {
		s = fmt.Sprint(v...)
	}
	p.Println(":", s)
}

func (p *Log) _fmtPrintln(a ...any) error {
	if p.mu != nil {
		p.mu.Lock()
		defer p.mu.Unlock()
	}
	return p.fmter.Println(a...)
}
func (p *Log) _fmtPrintf(format string, a ...any) error {
	if p.mu != nil {
		p.mu.Lock()
		defer p.mu.Unlock()
	}
	return p.fmter.Printf(format, a...)
}

type fmterDefault struct {
}

func (f *fmterDefault) Println(a ...any) error {
	_, err := fmt.Println(a...)
	return err
}
func (f *fmterDefault) Printf(format string, a ...any) error {
	_, err := fmt.Printf(format, a...)
	return err
}
