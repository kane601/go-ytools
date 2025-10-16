package tools

import (
	"sync"

	js "github.com/bitly/go-simplejson"
	"gopkg.in/ini.v1"
)

var (
	//ProductName 产品名
	_ProductName string
)

//ConfigValueFile ...
func ConfigValueFile(file string, keypath ...string) *js.Json {
	j, err := OpenJSON(file)
	if err != nil {
		return js.New()
	}
	if len(keypath) == 0 {
		return j
	}
	if len(keypath) == 1 && len(keypath[0]) == 0 {
		return j
	}
	return j.GetPath(keypath...)
}

func ConfigValueJson(j *js.Json, keypath ...string) *js.Json {
	if len(keypath) == 0 {
		return j
	}
	if len(keypath) == 1 && len(keypath[0]) == 0 {
		return j
	}
	return j.GetPath(keypath...)
}

var opendedJsonFileMutex sync.Mutex
var opendedJsonFile map[string]*js.Json = map[string]*js.Json{}

func openJsonFile(file string) *js.Json {
	opendedJsonFileMutex.Lock()
	defer opendedJsonFileMutex.Unlock()
	j, ok := opendedJsonFile[file]
	if ok {
		return j
	}
	j, err := OpenJSON(file)
	if err != nil {
		return js.New()
	}
	opendedJsonFile[file] = j
	return j
}

//ConfigValue 从config.json中获取
func ConfigValue(keypath ...string) *js.Json {
	j := openJsonFile(LocalPath("config.json"))
	return ConfigValueJson(j, keypath...)
}

//FrameConfigValue 从frameConfig.json中获取
func FrameConfigValue(keypath ...string) *js.Json {
	j := openJsonFile(LocalPath("frameConfig.json"))
	return ConfigValueJson(j, keypath...)
}

func SetProductName(product string) {
	_ProductName = product
}

//GetProductName ...
func GetProductName() (product string) {
	if len(_ProductName) == 0 {
		product = "unknowProduct"
		if j := FrameConfigValue("productName"); j != nil && len(j.MustString()) > 0 {
			product = j.MustString()
		}
		_ProductName = product
	} else {
		product = _ProductName
	}
	return
}

//ParseIniString ....
func ParseIniString(key, def, session, file string) string {
	cfg, err := ini.LoadSources(ini.LoadOptions{Insensitive: true}, file)
	if err != nil {
		return def
	}
	return cfg.Section(session).Key(key).String()
}
