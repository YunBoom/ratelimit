package ratelimit

import (
	"fmt"
	"github.com/yunboom/ratelimit/judge"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	rateLimit = 1 << 20 //限速1m/s
	src       = "/Users/zonst/Downloads/1.zip"
	dst       = "/Users/zonst/Downloads/2.zip"
)

var (
	readFile, writeFile *os.File
	err                 error
)

func init() {
	readFile, err = os.Open(src)
	if err != nil {
		panic(err)
	}
	writeFile, err = os.OpenFile(dst, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
}

func TestReadLimit(t *testing.T) {
	//限制读取速度
	now := time.Now()
	reader := NewStorage(rateLimit).Reader(readFile)
	if _, err = ioutil.ReadAll(reader); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("拷贝耗时:%v, 文件大小:%d", time.Since(now), reader.Size())
}

func TestWriteLimit(t *testing.T) {
	//限制写入速度
	now := time.Now()
	writer := NewStorage(rateLimit).Writer(writeFile)
	if _, err = io.Copy(readFile, writer); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("拷贝耗时:%v, 文件大小:%d", time.Since(now), writer.Size())
}

func TestCopyLimit(t *testing.T) {
	//限制复制速度
	now := time.Now()
	if size, err := NewStorage(rateLimit).Copy(readFile, writeFile); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("拷贝耗时:%v, 文件大小:%d", time.Since(now), size)
	}
}

func TestLimiter(t *testing.T) {
	//限制多少人限速
	judge.Add()
	defer judge.Done()
	now := time.Now()
	if judge.IsNeedLimit() {
		if size, err := NewStorage(rateLimit).Copy(readFile, writeFile); err != nil {
			t.Fatal(err)
		} else {
			fmt.Printf("拷贝耗时:%v, 文件大小:%d", time.Since(now), size)
		}
	}
}
