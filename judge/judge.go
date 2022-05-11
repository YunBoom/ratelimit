package judge

import (
	"sync/atomic"
	"time"
)

const (
	limitReq = 3
)

type Jude interface {
	IsNeedLimit() bool
}

var RateLimitJude = newSeqJude()

type SeqJude struct {
	seq    int32     //序号 1、2、3放行，之后的需要限速
	latest time.Time //上一次有人上传的时间
}

func newSeqJude() *SeqJude {
	return &SeqJude{
		seq:    0,
		latest: time.Now(),
	}
}

func (sl *SeqJude) IsNeedLimit() bool {
	//半个小时没人上传则不需要限速
	if time.Now().After(sl.latest.Add(time.Minute * 30)) {
		return false
	}

	var taken, isNeedLimit bool
	for !taken {
		seq := sl.GetSeq()
		if seq > limitReq {
			isNeedLimit = true
		}
		taken = atomic.CompareAndSwapInt32(&sl.seq, seq, seq)
	}

	return isNeedLimit
}

func (sl *SeqJude) Add() {
	sl.increment(1)
	sl.latest = time.Now() //不需要加锁，可以接受误差
}

func (sl *SeqJude) Done() {
	sl.increment(-1)
}

func (sl *SeqJude) increment(n int32) {
	atomic.AddInt32(&sl.seq, n)
}

func (sl *SeqJude) Clear() {
	atomic.StoreInt32(&sl.seq, 0)
}

func (sl *SeqJude) GetSeq() int32 {
	return atomic.LoadInt32(&sl.seq)
}

func GetSeq() int32 {
	return RateLimitJude.GetSeq()
}

func Add() {
	RateLimitJude.Add()
}

func Done() {
	RateLimitJude.Done()
}

func IsNeedLimit() bool {
	return RateLimitJude.IsNeedLimit()
}
