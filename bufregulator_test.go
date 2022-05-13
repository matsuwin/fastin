package fastin

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// New Handle
var fastIn = New(time.Second)

// 使用切片实现暂存容器
var bucket = make([]string, 0)

// 创建自动刷新异步任务，动态调节 bucket size 的大小
func init() {
	go func() {
		for {
			fastIn.Refresh()
		}
	}()
}

var mutex = &sync.Mutex{}

func write(data string) {

	// 并行条件下记得加锁
	mutex.Lock()
	defer mutex.Unlock()

	// 将数据放入暂存桶
	bucket = append(bucket, data)

	// 索引+1
	fastIn.Index++

	// 当桶满后进行数据刷盘
	wc := len(bucket)
	if wc >= fastIn.Size {
		bucketRefresh(wc)
	}
}

func bucketRefresh(wc int) {
	fmt.Printf("insert %d\n", wc)
	defer func() {
		bucket = nil
	}()

	// 批量写入磁盘或下游数据库
}

func Test(t *testing.T) {

	// 模拟大量数据写入
	for {
		time.Sleep(time.Millisecond * 10)

		// 写入一条数据
		go write("data")
	}
}
