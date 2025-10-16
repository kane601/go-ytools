package tools

import (
	"os"
	"runtime"
	"strings"
)

func GetPathEnvKey() string {
	keys := []string{"PATH", "path", "Path"}
	for _, k := range keys {
		path := os.Getenv(k)
		if len(path) > 0 {
			return k
		}
	}
	return "PATH"
}

func AppendPathEnv(p string) {
	k := GetPathEnvKey()
	path := os.Getenv(k)
	if runtime.GOOS == "windows" {
		path += ";"
	} else {
		path += ":"
	}
	path += p
	os.Setenv(k, path)
}

func IsWin() bool {
	return strings.EqualFold(runtime.GOOS, "windows")
}
