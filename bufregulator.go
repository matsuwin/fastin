package fastin

import "time"

/*
 * 缓冲区容量调节器
 *
 * 适用于大量数据的写入场景，实现分批次的落盘，减轻下游负载。
 */

type structure struct {
	Size  int // Size 当前缓冲区的容量
	Index int
	rate  time.Duration
}

func New(rate time.Duration) *structure {
	return &structure{rate: rate}
}

func (st *structure) Refresh() {
	defer func() {
		st.Size = st.Index
	}()
	st.Index = 1
	time.Sleep(st.rate)
}
