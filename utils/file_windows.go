package tools

import (
	"fmt"

	"github.com/kane601/go-w32/wutil"
)

// DiskUsage 获取路径的磁盘信息
func DiskUsage(path string) (disk wutil.DiskStatus, err error) {
	usag, b := wutil.DiskUsage(path)
	if !b {
		err = fmt.Errorf("Get DiskUsage Fail")
		return
	}
	disk = usag
	return
}

// GetFileVersion 获取文件版本信息
func GetFileVersion(file string) string {
	return wutil.GetFileVersion(file)
}
