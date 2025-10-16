package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ThePath 指定路径
func ThePath(root string, path ...string) string {
	root = AbsPath(root)
	tmp := append([]string{}, root)
	path = append(tmp, path...)
	return AbsJoinPath(path...)
}

func PathCutExt(path string) string {
	if path == "" {
		return ""
	}
	ext := filepath.Ext(path)
	if ext == "" {
		return path
	}
	p, f := cutSuffix(path, ext)
	if f {
		return p
	}
	return path
}

func cutSuffix(s, suffix string) (before string, found bool) {
	if !strings.HasSuffix(s, suffix) {
		return s, false
	}
	return s[:len(s)-len(suffix)], true
}

// AbsParent 绝对父路径
func AbsParent(path string) string {
	if len(path) == 0 {
		return ""
	}
	return filepath.Dir(AbsPath(path))
}

// AbsPath 绝对路径
func AbsPath(path string) string {
	return AbsJoinPath(path)
}

// AbsJoinPath 拼接路径
func AbsJoinPath(paths ...string) string {
	if 0 == len(paths) {
		return ""
	}
	abs, err := filepath.Abs(paths[0])
	if err != nil {
		if strings.HasPrefix(paths[0], "./") || strings.HasPrefix(paths[0], ".\\") {
			return "./" + filepath.Join(paths...)
		}
	}
	paths[0] = abs
	if strings.HasSuffix(paths[0], ":") {
		paths[0] += "/"
	}
	return filepath.Join(paths...)
}

// PathParent 父路径
func PathParent(path string) string {
	if len(path) == 0 || len(path) == 1 {
		return path
	}
	i := strings.LastIndex(path, "/")
	if i == -1 {
		i = strings.LastIndex(path, "\\")
	}
	if i == -1 {
		return ""
	}
	if i == 0 {
		return "/"
	}
	return path[0:i]
}

// JoinPath 拼接路径
func JoinPath(paths ...string) string {
	if len(paths) == 0 {
		return ""
	}
	if strings.HasPrefix(paths[0], "./") || strings.HasPrefix(paths[0], ".\\") {
		return "./" + filepath.Join(paths...)
	}
	if strings.HasSuffix(paths[0], ":") {
		paths[0] += "/"
	}
	return filepath.Join(paths...)
}

// PathName 从url或路径中获取文件名
func PathName(path string) string {
	index1 := strings.LastIndex(path, "/")
	index2 := strings.LastIndex(path, "\\")
	if -1 == index1 && -1 == index2 {
		return path
	}

	index := index1
	if index2 > index1 {
		index = index2
	}
	return path[index+1:]
}

// PathStem 从url或路径中获取不带dot的文件名
func PathStem(path string) string {
	dotname := PathName(path)
	index := strings.LastIndex(dotname, ".")
	if -1 != index {
		return dotname[:index]
	}
	return dotname
}

// PathAvailableSpace 可用空间
func PathAvailableSpace(path string) uint64 {
	usage, err := DiskUsage(path)
	if err != nil {
		return 0
	}
	return usage.Free
}

// ReplacePath 替换\为/
func ReplacePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// GetHasFileRoot 找到第一个非空路径
func GetHasFileRoot(root string) (ret string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.HasPrefix(info.Name(), ".") {
			return nil
		}
		if !info.IsDir() {
			ret = AbsParent(path)
			return fmt.Errorf("found")
		}
		return nil
	})
	return
}

// GetExtFilePath 从目录中获取指定后缀的文件路径
func GetExtFilePath(dirpath, ext string) string {
	ext = strings.ToLower(ext)
	dirpath = AbsPath(dirpath)
	rd, err := ioutil.ReadDir(dirpath)
	if nil != err {
		return ""
	}
	for _, fi := range rd {
		path := AbsJoinPath(dirpath, fi.Name())
		if fi.IsDir() {
			return GetExtFilePath(path, ext)
		} else if ext == strings.ToLower(filepath.Ext(path)) {
			return path
		}
	}
	return ""
}

// FilterFile 从目录中过滤文件,不会深度遍历
func FilterFile(dirpath string, filter []string) []string {
	dirpath = AbsPath(dirpath)
	rd, err := ioutil.ReadDir(dirpath)
	if nil != err {
		return []string{}
	}
	results := []string{}
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		if IsInFilter(fi.Name(), filter) {
			path := AbsJoinPath(dirpath, fi.Name())
			results = append(results, path)
		}
	}
	return results
}

// FilterDeepFile 从目录中过滤文件,会深度遍历
func FilterDeepFile(dirpath string, filter []string) []string {
	results := []string{}
	PathWalk(dirpath, func(path string, info os.FileInfo, postName string) error {
		if info.IsDir() {
			return nil
		}
		if IsInFilter(info.Name(), filter) {
			results = append(results, path)
		}
		return nil
	})
	return results
}

// IsInFilter 文件格式或名字是否在筛选器里面
func IsInFilter(file string, filter []string) bool {
	if len(filter) == 0 {
		return true
	}
	file = PathName(file)
	for _, pattern := range filter {
		b, _ := filepath.Match(pattern, file)
		if b {
			return b
		}
	}
	return false
}

// GetApplicationDir 获取当前程序目录
func GetApplicationDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

// IsEqualPath 是否是相同路径
func IsEqualPath(p1, p2 string) bool {
	if strings.HasSuffix(p1, ":") {
		p1 += "/"
	}
	if strings.HasSuffix(p2, ":") {
		p2 += "/"
	}

	p1 = AbsJoinPath(p1, "t")
	p2 = AbsJoinPath(p2, "t")

	p1 = ReplacePath(p1)
	p2 = ReplacePath(p2)

	p1 = strings.ToLower(p1)
	p2 = strings.ToLower(p2)
	return p1 == p2
}

// PathWalk 递归遍历目录，回调带上去除根路径的name
func PathWalk(path string, f func(p string, info os.FileInfo, postName string) error) {
	pathlen := len(path)
	filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if len(fpath) <= pathlen {
			return nil
		}
		if info.Name() == "." || info.Name() == ".." {
			return nil
		}
		postName := PostPath(fpath, path)
		return f(fpath, info, postName)
	})
}

// SpecialDirType 特定目录
type SpecialDirType int32

var (
	//LocalAppdata 数据路径
	LocalAppdata        SpecialDirType = 0x001c
	LocalRoamingAppdata SpecialDirType = 0x001a
	TempLocation        SpecialDirType = 0x1234
	PreferencesLocation SpecialDirType = 0x1235
	HomeLocation        SpecialDirType = 0x1236
)

var (
	_fixLocalPath   string
	_fixDataPath    string
	_fixRoamingPath string
)

// SetFixLocal ...
func SetFixLocal(path string) {
	_fixLocalPath = path
}

// SetFixDataPath ...
func SetFixDataPath(path string) {
	_fixDataPath = path
}

func SetFixRoamingDataPath(path string) {
	_fixRoamingPath = path
}

// LocalPath 当前程序文件路径
func LocalPath(path string) string {
	if len(_fixLocalPath) > 0 {
		return AbsJoinPath(_fixLocalPath, path)
	}
	return AbsJoinPath(AbsParent(os.Args[0]), path)
}

// PostPath full去除src的剩余路径路径
func PostPath(full, src string) string {
	full = strings.ReplaceAll(full, "\\", "/")
	src = strings.ReplaceAll(src, "\\", "/")

	i := strings.Index(full, src)
	if i == -1 {
		return ""
	}
	i = len(src)
	if full[i] == '\\' {
		i++
	}
	return full[i:]
}

// TempPath 临时目录文件路径
var specialTempDir string

func TempPath(path string) string {
	_specialDir := func() string {
		if specialTempDir != "" {
			return specialTempDir
		}
		specialTempDir = GetSpecialDir(TempLocation)
		return specialTempDir
	}

	return AbsJoinPath(splitPathProductName(_specialDir(), path)...)
}

// DataPath 数据目录文件路径
// :productname:path
var _specialDataDir string

func DataPath(path string) string {
	if len(_fixDataPath) > 0 {
		return AbsJoinPath(_fixDataPath, path)
	}
	_specialDir := func() string {
		if _specialDataDir != "" {
			return _specialDataDir
		}
		_specialDataDir = GetSpecialDir(LocalAppdata)
		return _specialDataDir
	}

	return AbsJoinPath(splitPathProductName(_specialDir(), path)...)
}

var _specialRoamingDataDir string

func RoamingDataPath(path string) string {
	if len(_fixRoamingPath) > 0 {
		return AbsJoinPath(_fixRoamingPath, path)
	}
	_specialDir := func() string {
		if _specialRoamingDataDir != "" {
			return _specialRoamingDataDir
		}
		_specialRoamingDataDir = GetSpecialDir(LocalRoamingAppdata)
		return _specialRoamingDataDir
	}

	return AbsJoinPath(splitPathProductName(_specialDir(), path)...)
}

// SetPathExt 改变路径中的后缀，不会修改文件
func SetPathExt(p, ext string) string {
	dir := filepath.Dir(p)
	name := PathStem(p)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	if !strings.HasSuffix(dir, "/") && !strings.HasSuffix(dir, "\\") {
		dir += "/"
	}
	return dir + name + ext
}

// PreferencesPath 数据目录文件路径
var specialPreferencesDir string

func PreferencesPath(path string) string {
	_specialDir := func() string {
		if specialPreferencesDir != "" {
			return specialPreferencesDir
		}
		specialPreferencesDir = GetSpecialDir(PreferencesLocation)
		return specialPreferencesDir
	}

	return AbsJoinPath(splitPathProductName(_specialDir(), path)...)
}

func GenUniquePath(p string) string {
	dir := AbsParent(p)
	name := PathStem(p)
	ext := filepath.Ext(p)
	for i := 1; IsExist(p); i++ {
		p = ThePath(dir, name+fmt.Sprintf("(%d)", i)+ext)
	}
	return p
}

func ReplaceIllegality(name string) string {
	if len(name) == 0 {
		return ""
	}
	re := regexp.MustCompile(`[\x{1}-\x{6}\x{e}-\x{19}\x{1b}-\x{1f}"<>\|\a\t\n\v\f\r\:\*\?\\\/]`)
	return re.ReplaceAllString(name, "_")
}

func CorreceLongPath(p string, reveried int) string {
	if len(p) <= reveried {
		return p
	}
	dir := AbsParent(p)
	name := PathStem(p)
	ext := filepath.Ext(p)

	nameLen := len(name) - (len(p) - reveried)
	if nameLen <= 0 {
		return p
	}
	if nameLen > 245 {
		nameLen = 245
	}
	p = ThePath(dir, name[0:nameLen]+ext)
	return p
}

func CorrectIllegalityPath(p string) string {
	name := ReplaceIllegality(PathStem(p))
	ext := filepath.Ext(p)
	parentDir := PathParent(p)
	finalDir := ""

	lastSpeIndex := -1

	_isSep := func(c rune) bool {
		return c == '\\' || c == '/'
	}

	for i := 0; i < len(parentDir); i++ {
		c := rune(parentDir[i])
		// 不是分隔符 也不是最后
		if !_isSep(c) && len(parentDir)-1 != i {
			continue
		}

		// 最后
		if len(parentDir)-1 == i {
			if lastSpeIndex != -1 {
				finalDir += "/" + ReplaceIllegality(parentDir[lastSpeIndex+1:])
			} else {
				finalDir = parentDir
			}
			break
		}

		// 根目录不修正
		if lastSpeIndex == -1 {
			finalDir = parentDir[:i]
			lastSpeIndex = i
			continue
		}

		//连续的分割符号
		if i-lastSpeIndex == 1 {
			finalDir += "/_"
			lastSpeIndex = i
			continue
		}

		// 分割
		finalDir += "/" + ReplaceIllegality(parentDir[lastSpeIndex+1:i])
		lastSpeIndex = i
	}

	path := JoinPath(finalDir, name+ext)
	path = CorreceLongPath(path, 260)
	return path
}

func IsWriteableDir(p string) bool {
	flagfile := AbsJoinPath(p, ".__is_writeable__")
	err := WriteFileString(flagfile, "explain this is writeable not delete")
	if err == nil {
		os.Remove(flagfile)
		return true
	}
	return false
}

func splitPathProductName(pre, path string) (ret []string) {
	productName := ""
	if len(path) > 0 && path[0] == '<' {
		i := strings.Index(path, ">")
		if i == -1 {
			productName = GetProductName()
		} else {
			productName = path[1:i]
			path = path[i+1:]
		}
	} else {
		productName = GetProductName()
	}

	ret = []string{}
	if len(pre) != 0 {
		ret = append(ret, pre)
	}
	if len(productName) != 0 {
		ret = append(ret, productName)
	}
	ret = append(ret, path)
	return
}
