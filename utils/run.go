package tools

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func CmdOutputs(name string, arg ...string) (stdoutLines, stderrLines []string) {
	var stdoutBuff, stderrBuff bytes.Buffer
	c := exec.Command(name, arg...)
	c.Stdout = &stdoutBuff
	c.Stderr = &stderrBuff
	c.Run()
	filterEmpty := func(arr []string) []string {
		res := []string{}
		for _, a := range arr {
			a = strings.TrimSpace(a)
			if a == "" {
				continue
			}
			res = append(res, a)
		}
		return res
	}
	stdoutLines = filterEmpty(strings.Split(strings.ReplaceAll(stdoutBuff.String(), "\r", "\n"), "\n"))
	stderrLines = filterEmpty(strings.Split(strings.ReplaceAll(stderrBuff.String(), "\r", "\n"), "\n"))
	return
}

func CmdRun(ctx context.Context,
	stdoutSink func(stdout io.Reader) error,
	stderrSink func(stderr io.Reader) error,
	is64 bool,
	wordDir string,
	name string, arg ...string) (stdoutSinkErr, stderrSinkErr, waitErr, startErr error) {

	cmd := exec.Command(name, arg...)
	if wordDir != "" {
		err := os.MkdirAll(string(wordDir), 0755)
		if err != nil {
			startErr = fmt.Errorf("MkdirAll(%s): %v", wordDir, err)
			return
		}
		cmd.Dir = string(wordDir)
	}

	var stderr io.ReadCloser
	if stderrSink != nil {
		stderr, startErr = cmd.StderrPipe()
		if startErr != nil {
			return
		}
		defer stderr.Close()
	}

	var stdout io.ReadCloser
	if stdoutSink != nil {
		stdout, startErr = cmd.StdoutPipe()
		if startErr != nil {
			return
		}
		defer stdout.Close()
	}

	startErr = cmd.Start()
	if startErr != nil {
		return
	}

	var cancleErr error
	c := make(chan struct{})
	defer close(c)
	go func() {
		select {
		case <-c:
			return
		case <-ctx.Done():
			if is64 {
				SendCtrlC64(cmd.Process)
			} else {
				SendCtrlC(cmd.Process)
			}
			cancleErr = ctx.Err()
		}
	}()

	wg := sync.WaitGroup{}

	if stderrSink != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stderrSinkErr = stderrSink(stderr)
		}()
	}

	if stdoutSink != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stdoutSinkErr = stdoutSink(stdout)
		}()
	}
	wg.Wait()
	waitErr = cmd.Wait()
	if cancleErr != nil {
		waitErr = cancleErr
	}
	return
}

type SplitTrimLineReader struct {
	r *bufio.Reader
}

func NewSplitTrimLineReader(r io.Reader) *SplitTrimLineReader {
	return &SplitTrimLineReader{
		r: bufio.NewReader(r),
	}
}

func (l *SplitTrimLineReader) Read() (line string, err error) {
	data := []byte{}
	for {
		var b byte
		b, err = l.r.ReadByte()
		if err != nil {
			return
		}
		if b == '\n' || b == '\r' {
			line = strings.TrimSpace(string(data))
			if line != "" {
				return
			}
		} else {
			data = append(data, b)
		}
	}
}
