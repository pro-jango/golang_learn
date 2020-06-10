#### 使用aotmic(原子操作)、锁机制保证协程更新数据的正确性
#### 问题： 
当多个协程更新同一个数据时，在不加锁的情况下会导致数据更新混乱。示例： 
```
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
//最终打印的结果不是1000
```

#### 使用锁机制保证正确性 
```
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
```
在每次对数据进行更新时，加上写锁，等协程操作完之后，再释放该写锁，这样其他程序将会阻塞，影响效率。

#### 采用原子操作实现 
```
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
```
我们和对比使用互斥锁实现等机制，发现通过原子操作的操作要比通过锁机制实现的快很多