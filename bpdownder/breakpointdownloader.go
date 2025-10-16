package bpd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	tools "github.com/kane601/go-ytools/utils"
)

// BPDownloader 断点下载器
type BPDownloader struct {
	url         string
	path        string
	offset      int64
	size        int64
	etagPath    string
	etagConeten string
	zeroTimeout int
}

// ProgHand ...
type ProgHand func(total, downed, speed int64, prog float64)

// New ...
func New() *BPDownloader {
	return &BPDownloader{}
}

// SetZeroTimeout ...
func (b *BPDownloader) SetZeroTimeout(timeout int) {
	b.zeroTimeout = timeout
}

// Download ...
func (b *BPDownloader) Download(ctx context.Context, url, path string, progFun ProgHand) (err error) {
	b.url = url
	b.path = path
	err = tools.CreateDirs(tools.AbsParent(path))
	if err != nil {
		return
	}

	err = b.fetchURLFileAttr()
	if err != nil {
		return
	}

	if b.offset == 0 {
		tools.RemovePath(b.etagPath)
		tools.RemovePath(b.path)
	} else if b.offset == b.size {
		return nil
	}
	b.downFile(ctx, progFun)
	return
}

func (b *BPDownloader) fetchURLFileAttr() (err error) {
	resp, err := http.Get(b.url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if !strings.EqualFold(resp.Header.Get("Accept-Ranges"), "bytes") {
		err = fmt.Errorf("No support Accept-Ranges")
		return
	}
	if resp.ContentLength <= 0 {
		err = fmt.Errorf("Not content length")
		return
	}

	etagConeten := resp.Header.Get("ETag") + resp.Header.Get("Last-modified")
	b.etagConeten, _ = tools.GenMd5([]byte(etagConeten))
	b.etagPath = b.path + ".yetag"

	if tools.IsExist(b.path) &&
		tools.IsExist(b.etagPath) {
		if etagConeten != tools.ReadFileString(b.etagPath) {
			fmt.Println("Last edittag is not equal")
		} else {
			b.offset = tools.FileSize(b.path)
		}
	}
	b.size = resp.ContentLength
	return
}

func (b *BPDownloader) downFile(ctx context.Context, progFun ProgHand) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if progFun == nil {
		progFun = func(total, downed, speed int64, prog float64) {}
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", b.url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", b.offset, b.size-1))
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := tools.OpenApptendFile(b.path)
	if err != nil {
		return
	}
	tools.WriteFileString(b.etagPath, b.etagConeten)

	finishChan := make(chan struct{}, 2)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case <-ctx.Done():
			resp.Body.Close()
		case <-finishChan:
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {

		last := b.offset
		zeroTime := 0
	loop:
		for {
			select {
			case <-finishChan:
				break loop
			case <-time.After(time.Second):
				if last != b.offset {
					progFun(b.size, b.offset, b.offset-last, float64(b.offset)/float64(b.size))
					last = b.offset
				} else {
					zeroTime++
					if b.zeroTimeout > 0 && zeroTime >= b.zeroTimeout {
						resp.Body.Close()
						break loop
					}
				}
			}
		}
		wg.Done()
	}()

	for b.size != b.offset {
		copyn := b.size - b.offset
		if copyn > 1024*1024 {
			copyn = 1024 * 1024
		}
		copyed, err := io.CopyN(file, resp.Body, copyn)
		b.offset += copyed

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	finishChan <- struct{}{}
	finishChan <- struct{}{}
	wg.Wait()
	return
}
