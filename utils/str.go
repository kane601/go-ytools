package tools

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unsafe"
)

//StringToFloat 字符串转浮点数
func StringToFloat(s string) (ret float64) {
	fmt.Sscanf(s, "%f", &ret)
	return
}

func StringToUnicode(s string) string {
	us := ""
	for _, a := range s {
		us += fmt.Sprintf(`\u%04x`, a)
	}
	return us
}

func UnicodeToString(us string) string {
	if !strings.HasPrefix(us, `\u`) {
		return us
	}

	_scan := func(b, e int) string {
		var c rune
		_, err := fmt.Sscanf(us[b+2:e], `%x`, &c)
		if err != nil {
			return ""
		}
		return string(c)
	}

	str := ""
	i := 0
	j := 2
	for ; j <= len(us); j++ {
		if j == len(us) || us[j] == '\\' {
			str += _scan(i, j)
			i = j
		}
	}
	return str
}

func StrMD5(url string) string {
	if len(url) == 0 {
		return ""
	}
	md5, err := GenMd5([]byte(url))
	if err != nil {
		md5 = "common_md5"
	}
	return md5
}

func StrSHA1(u string) string {
	hashInstance := crypto.SHA1.New()
	hashInstance.Write([]byte(u))
	hash := hex.EncodeToString(hashInstance.Sum(nil))
	return strings.ToUpper(hash)
}

func StringToDuration(dur string) (d time.Duration) {
	dur = strings.TrimSpace(dur)
	if len(dur) == 0 || len(dur) == 1 {
		return
	}
	durArr := strings.Split(dur, ":")
	unit := []string{"s", "m", "h"}
	durStr := ""
	j := 0
	for i := len(durArr) - 1; i >= 0 && j < 3; i-- {
		durStr = durArr[i] + unit[j] + durStr
		j++
	}
	d, _ = time.ParseDuration(durStr)
	return
}

func ByteToPtr(b []byte) uintptr {
	if len(b) == 0 {
		return uintptr(0)
	}
	return uintptr(unsafe.Pointer(&b[0]))
}
