package tools

import "syscall"

//DiskStatus ...
type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

//DiskUsage 获取路径的磁盘信息
func DiskUsage(path string) (disk DiskStatus, err error) {
	fs := syscall.Statfs_t{}
	err = syscall.Statfs(path, &fs)
	if err != nil {
		return
	}

	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

//GetFileVersion 只有windows版本
func GetFileVersion(file string) string {
	return ""
}
