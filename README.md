# ratelimit

#### 介绍
文件上传下载限速

#### 软件架构
软件架构说明


#### 安装教程

```go get github.com/yunboom/generate```

#### 使用说明

1. 读取限速
```go
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
```
2. 限制写入速度
```go
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

func TestWriteLimit(t *testing.T) {
//限制写入速度
now := time.Now()
writer := NewStorage(rateLimit).Writer(writeFile)
if _, err = io.Copy(readFile, writer); err != nil {
t.Fatal(err)
}
fmt.Printf("拷贝耗时:%v, 文件大小:%d", time.Since(now), writer.Size())
}
```

3. 限制复制速度
```go
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

func TestCopyLimit(t *testing.T) {
//限制复制速度
now := time.Now()
if size, err := NewStorage(rateLimit).Copy(readFile, writeFile); err != nil {
t.Fatal(err)
} else {
fmt.Printf("拷贝耗时:%v, 文件大小:%d", time.Since(now), size)
}
}
```

4.限制上传人数大于3时限速

```go
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
```

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

