# fast insert

> 试试快速写入

*缓冲区容量调节器，适用于大量数据的写入场景，实现分批次的落盘，减轻下游负载。*

<br>

## Quick Start

*1. 准备工作*

```go
// New Handle
var fastIn = fastin.New(time.Second)

// 使用切片实现暂存容器
var bucket = make([]string, 0)
```

*2. 创建刷新器*

```go
// 创建自动刷新异步任务，动态调节 bucket size 的大小
func init() {
    go func() {
        for {
            fastIn.Refresh(len(bucket))
        }
    }()
}
```

*3. 实现数据入桶和刷盘逻辑*

```go
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
    // ...
}
```

*4. 测试数据写入效果*

```go
// 模拟大量数据写入
for {
    time.Sleep(time.Millisecond * 10)
    
    // 写入一条数据
    go write("data")
}
```
