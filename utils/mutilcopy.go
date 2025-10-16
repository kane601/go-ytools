package tools

import (
	"context"
	"fmt"
	"io"
)

//ProgressHand 进度回调
type ProgressHand func(int64, float64)

//MutilCopyhander 多次拷贝处理
type MutilCopyhander struct {
	TotalSize int64
	Copyed    int64
	ProgHand  ProgressHand
	ctx       context.Context
}

//NewMutilCopyHander MutilCopyhander
func NewMutilCopyHander(ctx context.Context, totalSize int64, progHand ProgressHand) *MutilCopyhander {
	return &MutilCopyhander{
		TotalSize: totalSize,
		ProgHand:  progHand,
		Copyed:    0,
		ctx:       ctx,
	}
}

const copySize = 1024 * 100

//Copy 执行一次拷贝
func (c *MutilCopyhander) Copy(writer io.Writer, reader io.Reader) (timeWritenAll int64, err error) {
	if c.ProgHand == nil {
		c.ProgHand = func(int64, float64) {}
	}
	if c.TotalSize <= 0 {
		for {
			if c.ctx != nil {
				select {
				case <-c.ctx.Done():
					err = fmt.Errorf("Cancle Copy")
					return
				default:
				}
			}
			writen, err := io.CopyN(writer, reader, copySize)
			timeWritenAll += writen
			c.Copyed += writen
			c.ProgHand(c.TotalSize, float64(c.Copyed)/float64(c.TotalSize))

			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}
		}
		return
	}

	for c.TotalSize != c.Copyed {
		if c.ctx != nil {
			select {
			case <-c.ctx.Done():
				err = fmt.Errorf("Cancle Copy")
				return
			default:
			}
		}

		write := c.TotalSize - c.Copyed
		if write > copySize {
			write = copySize
		}
		writen, err := io.CopyN(writer, reader, write)
		timeWritenAll += writen
		c.Copyed += writen
		c.ProgHand(c.TotalSize, float64(c.Copyed)/float64(c.TotalSize))

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	return
}

//CopyFun 带回调的copy
func CopyFun(ctx context.Context, size int64, writer io.Writer, reader io.Reader, progf ProgressHand) (writen int64, err error) {
	mutilCopy := NewMutilCopyHander(ctx, size, progf)
	return mutilCopy.Copy(writer, reader)
}

//Copy ...
func Copy(ctx context.Context, writer io.Writer, reader io.Reader, progf ProgressHand) (writen int64, err error) {
	mutilCopy := NewMutilCopyHander(ctx, -1, progf)
	return mutilCopy.Copy(writer, reader)
}
