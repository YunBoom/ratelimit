package judge

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestSeqJude_Increment(t *testing.T) {
	jude := newSeqJude()
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			jude.Add()
			wg.Done()
		}()
	}

	wg.Wait()
	if jude.GetSeq() != 10000 {
		t.Errorf("want 10000 but get %d", jude.GetSeq())
	}
}

func TestSeqJude_Clear(t *testing.T) {
	jude := newSeqJude()
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(j int) {
			jude.Add()
			wg.Done()
		}(i)
	}

	wg.Wait()
	jude.Clear()
	if jude.GetSeq() != 0 {
		t.Errorf("want 0 but get %d", jude.GetSeq())
	}
}

func TestSeqJude_IsNeedLimit(t *testing.T) {
	var n int32
	jude := newSeqJude()
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			jude.Add()
			if !jude.IsNeedLimit() {
				atomic.AddInt32(&n, 1)
				return
			}
		}()
	}

	wg.Wait()

	if n != 3 {
		t.Errorf("want 3 but get %d", n)
	}

	fmt.Println(jude.GetSeq())
}

func TestSeqJude_Done(t *testing.T) {
	jude := newSeqJude()
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			jude.Add()
		}()
	}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			jude.Done()
		}()
	}

	wg.Wait()

	if jude.GetSeq() != 0 {
		t.Errorf("want 0 but get %d", jude.GetSeq())
	}
}
