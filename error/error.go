package yerror

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	_ Error        = (*_Error)(nil)
	_ fmt.Stringer = (*_Error)(nil)
)

//Error ...
type Error interface {
	error
	fmt.Stringer
	Code() int
	Caller() []CallerInfo
	CallerInfoStr() string
	private()
}

type _Error struct {
	err     error
	code    int
	callers []CallerInfo
	wraped  error
}

var _codeTranslate = map[int]string{}

//AddCodeTranslateMap 添加错误码映射关系
func AddCodeTranslateMap(codeMsg map[int]string) {
	for code, msg := range codeMsg {
		_codeTranslate[code] = msg
	}
}

//GetCodeTranslate 获取错误码对应的翻译
func GetCodeTranslate(code int) string {
	msg, ok := _codeTranslate[code]
	if !ok {
		if code == 0 {
			return "Successed"
		}
		return ""
	}
	return msg
}

func ErrorSplit(err error) (code int, str string, valid bool) {
	if e, ok := err.(Error); ok {
		code = e.Code()
		str = e.String()
		valid = true
	}
	return
}

//New ...
func New(e interface{}, v ...interface{}) Error {
	if e == nil {
		return nil
	}
	code := -1
	var endmsg string
	if c, ok := e.(int); ok {
		code = c
		if code == 0 {
			return nil
		}
	} else if yerr, ok := e.(Error); ok {
		code = yerr.Code()
		endmsg = yerr.Error()
	} else if err, ok := e.(error); ok {
		code = -1
		endmsg = err.Error()
	} else {
		code = -1
		endmsg = fmt.Sprint(e)
	}

	var msg string
	if len(v) > 0 {
		if fmat, ok := v[0].(string); ok && strings.HasPrefix(fmat, "format!") {
			fmat := fmat[7:]
			v = v[1:]
			msg = fmt.Sprintf(fmat, v...)
		} else {
			msg = fmt.Sprint(v...)
		}
	}
	if len(msg) > 0 && len(endmsg) > 0 {
		rmsg := []rune(msg)
		rend := []rune(endmsg)
		i := rmsg[len(rmsg)-1]
		j := rend[0]
		if i != '.' && i != '。' && i != '，' && i != ',' &&
			j != '.' && j != '。' && j != '，' && j != ',' {
			msg += "," + endmsg
		} else {
			msg += endmsg
		}

	} else if len(endmsg) > 0 {
		msg = endmsg
	}
	return &_Error{
		code:    code,
		err:     errors.New(msg),
		callers: CallerList(1),
	}
}

func (p *_Error) Code() int {
	return p.code
}

func (p *_Error) Caller() []CallerInfo {
	return p.callers
}

func (p *_Error) CallerInfoStr() (ret string) {
	first := true
	i := 0
	for _, info := range p.callers {
		if first {
			first = false
		} else {
			ret += " => "
		}
		ret += "file:" + info.File + ",func:" + info.Fun + ",line:" + strconv.Itoa(info.Line)
		i++
		if i == 3 {
			break
		}
	}
	return
}

func (p *_Error) Error() string {
	if nil == p.err {
		return GetCodeTranslate(p.code)
	}
	return p.err.Error()
}

func (p *_Error) String() (ret string) {
	ret = p.Error()
	return
}

func (p *_Error) private() {
	panic("unreached")
}
