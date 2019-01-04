package id

import (
	"sync"
	"time"
)

const (
	nodeBits  uint8 = 10
	stepBits  uint8 = 12
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	stepMax   int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

// 起始时间戳 (毫秒数显示)
var epoch int64 = 1514764800000 //2018-01-01 00:00:00

// 存储基础信息的 snowflake 结构
type snowflake struct {
	mu        sync.Mutex // 保证并发安全
	timestamp int64
	node      int64
	step      int64
}

// 生成、返回唯一 snowflake ID
func (n *snowflake) Generate() int64 {
	n.mu.Lock()         // 保证并发安全, 加锁
	defer n.mu.Unlock() // 解锁

	// 获取当前时间的时间戳 (毫秒数显示)
	now := time.Now().UnixNano() / 1e6

	if n.timestamp == now {
		// step 步进 1
		n.step++

		// 当前 step 用完
		if n.step > stepMax {
			// 等待本毫秒结束
			for now <= n.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 本毫秒内 step 用完
		n.step = 0
	}

	n.timestamp = now

	result := (now-epoch)<<timeShift | (n.node << nodeShift) | (n.step)

	return result
}
