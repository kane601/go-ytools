package tools

import (
	"unsafe"
)

//StringToWCharPtr mac
func StringToWCharPtr(str string, pstr *[]byte) {
	if len(str) == 0 {
		*pstr = []byte{0, 0, 0, 0}
		return
	}

	*pstr = append([]byte{}, (*[1 << 29]byte)(unsafe.Pointer(&([]rune(str)[0])))[0:len(str)*4]...)
	*pstr = append(*pstr, 0, 0, 0, 0)
}

//WCharByteToString mac
func WCharByteToString(buff []byte) string {
	if len(buff) == 0 {
		return ""
	}
	u32 := (*[1 << 29]rune)(unsafe.Pointer(&buff[0]))[0 : len(buff)/4 : len(buff)/4]
	return string(u32)
}
