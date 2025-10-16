package tools

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// DownFile 下载文件
func DownFile(ctx context.Context, url, path string) error {
	return DownFileFun(ctx, url, path, nil)
}

// DownFileFun 下载文件
func DownFileFun(ctx context.Context, url, path string, progHand ProgressHand) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Get url fail,url:%s,err:%v", url, err)
	}

	var fsize int64
	if len(resp.Header.Values("Content-Length")) > 0 {
		fsize, _ = strconv.ParseInt(resp.Header.Values("Content-Length")[0], 10, 64)
	}

	defer resp.Body.Close()
	file, err := CreateFile(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = CopyFun(ctx, fsize, file, resp.Body, progHand)
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("Copy body fail,url:%s,err:%v", url, err)
	}
	return nil
}

func GetURLFileSize(url string) int64 {
	response, err := http.Head(url)
	if err != nil {
		return 0
	}
	defer response.Body.Close()

	fileSize := response.Header.Get("Content-Length")

	fileSizeInt, err := strconv.ParseInt(fileSize, 10, 64)
	if err != nil {
		return 0
	}
	return fileSizeInt
}

// IsValidURL  url是否有效
func IsValidURL(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return 200 == res.StatusCode
}

// URLHeaderAttr url返回的数据头的属性值
func URLHeaderAttr(url, name string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ``
	}
	return resp.Header.Values(name)[0]
}

// GetNetContent ...
func GetNetContent(url string) (ret []byte, e error) {
	client := HttpClient()
	resp, err := client.Get(url)
	if err != nil {
		e = err
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		e = err
		return
	}
	ret = buf.Bytes()

	return
}

// GetNetContent ...
func GetNetContentAndCode(url string) (ret []byte, code int, e error) {
	client := HttpClient()

	resp, err := client.Get(url)
	if err != nil {
		e = err
		return
	}
	code = resp.StatusCode
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		e = err
		return
	}
	ret = buf.Bytes()
	return
}

func IsNetErrorCode(code int) bool {
	return code >= 400 && code <= 504
}

// GetNetCode ...
func GetNetCode(url string) (code int, e error) {
	client := HttpClient()

	resp, err := client.Get(url)
	if err != nil {
		e = err
		return
	}
	code = resp.StatusCode
	return
}

// SendNetRequest ...
func SendNetRequest(method, url string, head map[string]string, body io.Reader, recv io.Writer) (http.Header, error) {
	client := HttpClient()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		err = fmt.Errorf("SendNetRequest fail,%v", err)
		return nil, err
	}
	for key, val := range head {
		req.Header.Add(key, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("SendNetRequest fail: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if recv != nil {
		_, err = io.Copy(recv, resp.Body)
		if err != nil && err != io.EOF {
			err = fmt.Errorf("SendNetRequest fail,%v", err)
			return nil, err
		}
	}
	return resp.Header.Clone(), nil
}

var Proxy string

func SetProxy(proxy string) {
	Proxy = proxy
}

func HttpClient() *http.Client {
	if Proxy == "" {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 60,
		}
	}
	proxyUrl, err := url.Parse(Proxy)
	if err != nil {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 60,
		}
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 60,
	}
}
