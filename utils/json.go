package tools

import (
	"io"
	"os"
	"reflect"
	"strings"

	js "github.com/bitly/go-simplejson"
)

func isContain(s []interface{}, f interface{}) bool {
	for _, val := range s {
		if reflect.DeepEqual(val, f) {
			return true
		}
	}
	return false
}

func mergeJSONArray(toarray []interface{}, fromarray []interface{}) []interface{} {
	for _, fval := range fromarray {
		if !isContain(toarray, fval) {
			toarray = append(toarray, fval)
		}
	}
	return toarray
}

func mergeJSONMap(tomap map[string]interface{}, frommap map[string]interface{}) {
	for fkey, fval := range frommap {
		tval, ok := tomap[fkey]
		if !ok {
			tomap[fkey] = fval
			continue
		}

		typeof := reflect.TypeOf(fval)
		switch typeof.Kind() {
		case reflect.Slice:
			tomap[fkey] = mergeJSONArray(tval.([]interface{}), fval.([]interface{}))
		case reflect.Map:
			mergeJSONMap(tval.(map[string]interface{}), fval.(map[string]interface{}))
		default:
			tomap[fkey] = fval
		}
	}
}

//MergeJSON from合并json到to,以from为准
func MergeJSON(to *js.Json, from *js.Json) {
	totype := reflect.TypeOf(to.Interface())
	fromtype := reflect.TypeOf(from.Interface())
	if totype.Kind() != fromtype.Kind() {
		return
	}

	switch totype.Kind() {
	case reflect.Slice:
		mergeJSONArray(to.MustArray(), from.MustArray())
	case reflect.Map:
		mergeJSONMap(to.MustMap(), from.MustMap())
	}

}

//MergeFileJSON from合并json到to,以from为准
func MergeFileJSON(to io.Reader, from io.Reader, result io.Writer) error {
	fromjson, err := js.NewFromReader(from)
	if err != nil {
		return err
	}
	tojson, err := js.NewFromReader(to)
	if err != nil {
		return err
	}

	MergeJSON(tojson, fromjson)
	data, err := tojson.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = result.Write(data)
	if err != nil {
		return err
	}
	return nil
}

//OpenJSON ...
func OpenJSON(path string) (*js.Json, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return js.NewFromReader(file)
}

//CreateJSON 创建json节点
func CreateJSON() *js.Json {
	return js.New()
}

//MarshalToJSON ...
func MarshalToJSON(path string, j *js.Json) error {
	data, err := j.MarshalJSON()
	if err != nil {
		return err
	}
	return WriteFileString(path, string(data))
}

func FoundJSONString(data string, obj bool) (ret string) {
	startKey := '{'
	endKey := '}'
	if !obj {
		startKey = '['
		endKey = ']'
	}

	s := strings.Index(data, string(startKey))
	if s == -1 {
		return
	}
	data = data[s:]
	num := 0
	for i, c := range data {
		if c == rune(startKey) {
			num++
		} else if c == rune(endKey) {
			num--
		}
		if num == 0 {
			ret = data[:i+1]
			break
		}
	}
	return
}
