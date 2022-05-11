package ratelimit

import (
	"bytes"
	"github.com/fwhezfwhez/errorx"
	"io"
	"os"
	"strings"
)

const (
	singleCopySize = 2 * 1 << 10 //单次 copy 2kb
)

type Storage struct {
	limiter *fileLimiter
	reader  io.Reader
	writer  io.Writer
}

func NewStorage(rate int64) *Storage {
	return &Storage{
		limiter: NewFileLimiter(rate),
	}
}

func (rate *Storage) Reader(reader io.Reader) *Storage {
	rate.reader = reader
	return rate
}

func (rate *Storage) Writer(writer io.Writer) *Storage {
	rate.writer = writer
	return rate
}

func (rate *Storage) Read(buf []byte) (n int, err error) {
	if rate.reader == nil {
		return 0, errorx.NewFromString("reader is nil")
	}
	if n, err = rate.reader.Read(buf); err != nil {
		return 0, err
	}

	rate.limiter.Wait(int64(len(buf)))
	return n, nil
}

func (rate *Storage) Write(buf []byte) (n int, err error) {
	if rate.writer == nil {
		return 0, errorx.NewFromString("writer is nil")
	}

	if n, err = rate.writer.Write(buf); err != nil {
		return 0, err
	}

	rate.limiter.Wait(int64(len(buf)))
	return n, nil
}

// Size 存储桶计算content-length需要
func (rate Storage) Size() (length int64) {
	if rate.reader == nil {
		return 0
	}

	switch v := rate.reader.(type) {
	case *bytes.Buffer:
		return int64(v.Len())
	case *bytes.Reader:
		return int64(v.Len())
	case *strings.Reader:
		return int64(v.Len())
	case *os.File:
		stat, err := v.Stat()
		if err != nil {
			return 0
		}
		return stat.Size()
	case *io.LimitedReader:
		return v.N
	default:
		return 0
	}
}

func (rate *Storage) Copy(src io.Reader, dst io.Writer) (int64, error) {
	if src == nil || dst == nil {
		return 0, errorx.NewFromString("src or dst is nil")
	}
	rate.limiter.mode = DualChannel
	defer func() { rate.limiter.mode = SingleChannel }()
	rate.reader = src
	rate.writer = dst
	buf := make([]byte, singleCopySize)

	written, err := io.CopyBuffer(rate, rate, buf)
	return written, errorx.Wrap(err)
}
