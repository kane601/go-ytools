package tools

import (
	"encoding/json"
	"io"
)

// MarshalToJSONFile 序列化到文件
func MarshalToJSONFile(s interface{}, path string) error {
	file, err := CreateFile(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return MarshalToJSONWriter(s, file)
}

func MarshalToJSONWriter(s interface{}, w io.Writer) error {
	jsonMarsh1, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonMarsh1)
	return err
}

// UnMarshalJSONFile 反序列化JSON
func UnMarshalJSONFile(path string, out interface{}) error {
	file, err := OpenReadFile(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return UnMarshalJSONReader(file, out)
}

func UnMarshalJSONReader(r io.Reader, out interface{}) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return err
	}
	return nil
}
