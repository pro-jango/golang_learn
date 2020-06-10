package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

func main() {
	numAddTest()
	lockAddTest()
	aotmicAddTest()
}

func numAddTest() {
	var total = 0
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			total = total + 1
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("func:%s\n", "numAddTest")
	fmt.Printf("total:%d\n\n", total)
}

func lockAddTest() {
	var total = 0
	var l sync.RWMutex
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			l.Lock()
			defer l.Unlock()
			total = total + 1
			time.Sleep(time.Microsecond * 500)
			wg.Done()
		}()
	}
	wg.Wait()
	elapsed := (int)(time.Since(start) / time.Millisecond)
	fmt.Printf("func:%s\n", "lockAddTest")
	fmt.Printf("time consuming:%d\n", elapsed)
	fmt.Printf("total:%d\n\n", total)
}

//aotmicAddTest 利用aotmic的原子操作实现加法
func aotmicAddTest() {
	var total AotmicInt
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			total.Add(1)
			time.Sleep(time.Microsecond * 500)
			wg.Done()
		}()
	}
	wg.Wait()

	elapsed := (int)(time.Since(start) / time.Millisecond)
	fmt.Printf("func:%s\n", "aotmicAddTest")
	fmt.Printf("time consuming:%d\n", elapsed)
	fmt.Printf("total:%d\n\n", total)
}

type AotmicInt int64

func (a *AotmicInt) Add(val int64) {
	for {
		//指针类型强制转换需要借助unsafe.Pointer完成
		b := (*int64)(unsafe.Pointer(a))
		oldVal := atomic.LoadInt64(b)
		newVal := *b + val
		if atomic.CompareAndSwapInt64(b, oldVal, newVal) {
			return
		}
	}
}
