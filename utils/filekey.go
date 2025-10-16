package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	_filekeyDefRoot string
)

type FileKey string

//SetDefFileKeyRoot ...
func SetDefFileKeyRoot(r string) {
	_filekeyDefRoot = r
}

//GetFileKeyDefRoot ...
func GetFileKeyDefRoot() string {
	if len(_filekeyDefRoot) > 0 {
		return _filekeyDefRoot
	}
	return TempPath("filekey")
}

//IsFileKey xxxxxx&ext
func IsFileKey(s string) bool {
	index := strings.Index(s, "&")
	if index == -1 {
		return false
	}
	i, err := strconv.ParseUint("0x"+s[:index], 0, 64)
	if err != nil {
		return false
	}
	return i != 0
}

func ClearFileKeyRoot() {
	RemovePath(GetFileKeyDefRoot())
}

//NewFileKey 如果文件已存在，则会move
func NewFileKey(p string, storeOriName bool) FileKey {
	name := PathName(p)
	ext := filepath.Ext(name)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	var n int64
	for n == 0 || n == -1 {
		n = time.Now().UnixNano()<<32 | (RandNumN(32) + int64(len(name)))
	}

	filekey := FileKey(strconv.FormatUint(uint64(n), 16) + "&" + ext)
	if storeOriName {
		filekey.SetOriName(name)
	}

	if IsExist(p) {
		MoveFile(p, filekey.MapPath())
	}
	return filekey
}

func (k FileKey) Ext() string {
	_, ext := k.splitFileKey()
	return ext
}

func (k FileKey) MapPath() string {
	name, ext := k.splitFileKey()
	if len(ext) == 0 {
		return ThePath(GetFileKeyDefRoot(), name)
	}
	return ThePath(GetFileKeyDefRoot(), name+"."+ext)
}

func (k FileKey) Dir() string {
	return GetFileKeyDefRoot()
}

func (k FileKey) Open() (f *os.File, err error) {
	return OpenReadFile(k.MapPath())
}

func (k FileKey) Delete() {
	RemovePath(k.MapPath())
}

func (k FileKey) Create() (f *os.File, err error) {
	return CreateFile(k.MapPath())
}

func (k FileKey) GetOriStem() (stem string, err error) {
	name, err := k.GetOriName()
	if err != nil {
		return
	}
	stem = PathStem(name)
	return
}

func (k FileKey) GetOriName() (name string, err error) {
	val, err := k.GetStoredValue("name")
	if err != nil {
		return
	}
	name, ok := val.(string)
	if !ok {
		err = fmt.Errorf("key val is not a string")
		return
	}
	return
}

func (k FileKey) SetOriName(name string) error {
	return k.StoreKeyValue("name", name)
}

func (k FileKey) StoreKeyValue(key string, val interface{}) error {
	info, err := k.GetInfo()
	if err != nil {
		info = make(map[string]interface{})
		err = nil
	}
	info[key] = val
	return k.StoreInfo(info)
}

func (k FileKey) GetStoredValue(key string) (val interface{}, err error) {
	info, err := k.GetInfo()
	if err != nil {
		return
	}
	val, ok := info[key]
	if !ok {
		err = fmt.Errorf("Not found key")
		return
	}
	return
}

func (k FileKey) GetInfo() (info map[string]interface{}, err error) {
	err = UnMarshalJSONFile(k.MapPath()+".info", &info)
	return
}

func (k FileKey) AddStoreInfo(info map[string]interface{}) error {
	sinfo, err := k.GetInfo()
	if err != nil {
		sinfo = make(map[string]interface{})
		err = nil
	}
	for k, i := range info {
		sinfo[k] = i
	}
	return k.StoreInfo(sinfo)
}

func (k FileKey) StoreInfo(info map[string]interface{}) error {
	return MarshalToJSONFile(info, k.MapPath()+".info")
}

func (k FileKey) splitFileKey() (name, ext string) {
	ss := strings.Split(string(k), "&")
	if len(ss) != 2 && len(ss) != 1 {
		return
	}
	name = ss[0]
	if len(ss) == 2 {
		ext = ss[1]
	}
	return
}
