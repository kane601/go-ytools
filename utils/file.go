package tools

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// IsExist 文件是否存在
func IsExist(path string) bool {
	if 0 == len(path) {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func AddExeSuff(path string) string {
	if IsWin() {
		if strings.HasSuffix(strings.ToLower(path), ".exe") {
			return path
		}
		return path + ".exe"
	}
	return path
}

// MoveFile 移动文件
func MoveFile(src, dst string) error {
	if err := CreateDirs(filepath.Dir(dst)); err != nil {
		return err
	}
	os.Remove(dst)
	return os.Rename(src, dst)
}

// MoveFileTo 移动文件
func MoveFileTo(src, dir string) error {
	if err := CreateDirs(dir); err != nil {
		return err
	}
	return os.Rename(src, AbsJoinPath(dir, filepath.Base(src)))
}

func SetFileExt(p, ext string) error {
	newp := SetPathExt(p, ext)
	return os.Rename(p, newp)
}

// CreateFile 创建打开文件
func CreateFile(path string, fileperm ...fs.FileMode) (*os.File, error) {
	path = AbsPath(path)
	CreateDirs(filepath.Dir(path))

	if len(fileperm) == 0 {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		return file, err
	}
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileperm[0])
}

func CreateUniqueFile(p string) (*os.File, error) {
	return CreateFile(GenUniquePath(p))
}

// CopyFile 拷贝文件
func CopyFile(ctx context.Context, src string, dst string) error {
	return CopyFileFun(ctx, src, dst, nil)
}

// CopyFileFun 拷贝文件带回调
func CopyFileFun(ctx context.Context, src string, dst string, progf func(int64, float64)) error {
	size := FileSize(src)
	filesrc, err := OpenReadFile(src)
	if err != nil {
		return err
	}
	defer filesrc.Close()
	filedst, err := CreateFile(dst)
	if err != nil {
		return err
	}
	defer filedst.Close()
	_, err = CopyFun(ctx, size, filedst, filesrc, progf)
	return err
}

// RenameFile 重命名文件
func RenameFile(path, name string) string {
	if PathName(path) == name {
		return path
	}
	path = AbsPath(path)
	newpath := filepath.Join(filepath.Dir(path), name)
	RemovePath(newpath)
	os.Rename(path, newpath)
	return newpath
}

func RenameFileSteam(path, steam string) string {
	if PathStem(path) == steam {
		return path
	}
	ext := filepath.Ext(path)
	newpath := filepath.Join(filepath.Dir(path), steam+ext)
	RemovePath(newpath)
	os.Rename(path, newpath)
	return newpath
}

// OpenReadFile 打开读取文件
func OpenReadFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, err
}

// OpenApptendFile 追加方式打开或创建文件
func OpenApptendFile(path string) (*os.File, error) {
	CreateDirs(AbsParent(path))
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_RDONLY|os.O_CREATE, 0644)
}

// FileSize 文件大小
func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// IsSameFile 是否相同的文件
func IsSameFile(path1, path2 string) bool {
	m1, err := GenFileMd5(path1)
	if err != nil {
		return false
	}
	m2, err := GenFileMd5(path2)
	if err != nil {
		return false
	}
	if 0 == len(m1) && 0 == len(m2) {
		return false
	}
	return m1 == m2
}

// CreateDirs 创建目录
func CreateDirs(path string, dirperm ...fs.FileMode) error {
	path = AbsPath(path)
	if 1 == len(path) {
		return nil
	}
	var perm fs.FileMode = 0755
	if len(dirperm) > 0 && dirperm[0] != 0 {
		perm = dirperm[0]
	}
	err := os.MkdirAll(path, perm)
	return err
}

// ReadFileAll 读取整个文件
func ReadFileAll(path string) ([]byte, error) {
	file, err := OpenReadFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RemovePath 删除路径
func RemovePath(path string) {
	st, err := os.Stat(path)
	if err != nil {
		return
	}
	if st.IsDir() {
		os.RemoveAll(path)
	}
	os.Remove(path)
}

// RemoveFile 删除文件,如果是目录则忽略
func RemoveFile(path string) {
	st, err := os.Stat(path)
	if err != nil {
		return
	}
	if st.IsDir() {
		return
	}
	os.Remove(path)
}

// WriteFileString 向文件中写入字符串
func WriteFileString(path string, content string) error {
	file, err := CreateFile(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

// ReadFileString 读取文件内容字符串
func ReadFileString(path string) string {
	content, _ := ReadFileAll(path)
	if len(content) > 0 {
		return string(content)
	}
	return ""
}

// FileModifyTime 返回UTC的文件修改时间
func FileModifyTime(path string) (ret int64) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	ret = info.ModTime().UTC().Unix()
	return
}

// CheckDirPermission 检查目录权限
func CheckDirPermission(dir string) bool {
	if !IsExist(dir) {
		err := CreateDirs(dir)
		if err != nil {
			return false
		}
		os.Remove(dir)
		return true
	}
	file, err := CreateFile(AbsJoinPath(dir, "_check_dir_primission.tmp"))
	if err != nil {
		return false
	}
	file.Close()
	RemovePath(AbsJoinPath(dir, "_check_dir_primission.tmp"))
	return true
}

// CheckReadPermission 检查文件读权限
func CheckReadPermission(path string) bool {
	file, err := OpenReadFile(path)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

// IsUndamagedFile 是否是未损坏的文件
func IsUndamagedFile(path, md5 string) bool {
	if !IsExist(path) {
		return false
	}
	fmd5, _ := GenFileMd5(path)
	return fmd5 == md5
}

// IsDir 是否是目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// PathSize 指定路径大小
func PathSize(path string) (size int64) {
	if !IsDir(path) {
		return FileSize(path)
	}

	PathWalk(path, func(p string, info os.FileInfo, postName string) error {
		if !info.IsDir() {
			size += FileSize(p)
		}
		return nil
	})
	return
}

// CopyDirFiles 递归拷贝目录中的文件到指定目录
func CopyDirFiles(ctx context.Context, src, dst string) {
	PathWalk(src, func(path string, info os.FileInfo, postName string) error {
		if !info.IsDir() {
			CopyFile(ctx, path, AbsJoinPath(dst, postName))
		} else {
			CreateDirs(AbsJoinPath(dst, postName))
		}
		if IsContextDone(ctx) {
			return fmt.Errorf("finish")
		}
		return nil
	})
}

func CopyGlobFiles(ctx context.Context, glob, dst string) {
	if IsDir(glob) {
		CopyDirFiles(ctx, glob, ThePath(dst, PathName(glob)))
		return
	}

	matchs, err := filepath.Glob(glob)
	if err != nil {
		return
	}

	for _, s := range matchs {
		if IsContextDone(ctx) {
			break
		}
		d := ThePath(dst, PathName(s))
		if !IsDir(s) {
			CopyFile(ctx, s, d)
		} else {
			CreateDirs(s)
		}
	}
}

func CopyGlobFeepFiles(ctx context.Context, glob, dst string) {
	if IsDir(glob) {
		CopyDirFiles(ctx, glob, ThePath(dst, PathName(glob)))
		return
	}

	matchs, err := filepath.Glob(glob)
	if err != nil {
		return
	}

	for _, s := range matchs {
		if IsContextDone(ctx) {
			break
		}
		d := ThePath(dst, PathName(s))
		if !IsDir(s) {
			CopyFile(ctx, s, d)
		} else {
			CopyDirFiles(ctx, s, d)
		}
	}
}

// ReplaceFileTo 根据map替换文本内容到
func ReplaceFileTo(path string, m map[string]string, dst string) error {
	data, err := ReadFileAll(path)
	if err != nil {
		return err
	}
	strdata := ReplaceString(string(data), m)
	return WriteFileString(dst, strdata)
}
