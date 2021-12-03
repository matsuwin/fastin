package bufregulator

import (
	"fmt"
	"testing"
	"time"
)

// New bufregulator handle
var bufRegulator = New(time.Second)

// 业务数据暂存容器切片
var bucket = make([]string, 0)

func Test(t *testing.T) {

	// 创建自动刷新任务，动态调节 bucket size 的大小
	go func() {
		for {
			bufRegulator.Refresh(len(bucket))
		}
	}()

	// 模拟大量数据写入
	for {
		time.Sleep(time.Millisecond * 10)

		// 写入一条数据
		write("data")
	}
}

func write(data string) {

	// 将数据放入暂存桶 (动态分桶)
	bucket = append(bucket, data)

	// 索引 +1
	bufRegulator.Index++

	// 当桶满后进行数据刷盘
	wc := len(bucket)
	if wc >= bufRegulator.Size {
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
