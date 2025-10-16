package tools

import (
	"syscall"
	"unsafe"
)

//StringToWCharPtr windows
func StringToWCharPtr(str string, pstr *[]byte) {
	if len(str) == 0 {
		*pstr = []byte{0, 0}
		return
	}
	p, _ := syscall.UTF16PtrFromString(str)

	*pstr = append([]byte{}, (*[1 << 29]byte)(unsafe.Pointer(p))[0:len(str)*2]...)
	*pstr = append(*pstr, 0, 0)
}

//WCharByteToString windows
func WCharByteToString(buff []byte) string {
	if len(buff) == 0 {
		return ""
	}
	u16 := (*[1 << 29]uint16)(unsafe.Pointer(&buff[0]))[0 : len(buff)/2 : len(buff)/2]
	return syscall.UTF16ToString(u16)
}
