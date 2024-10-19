package id

import (
	"sync"
	"time"
)

// Snowflake 结构体
type Snowflake struct {
	mutex          sync.Mutex // 互斥锁
	epoch          int64      // 起始时间戳（毫秒）
	nodeIDBits     uint       // 节点ID所占的位数
	sequenceBits   uint       // 序列号所占的位数
	nodeIDShift    uint       // 节点左移位数
	timestampShift uint       // 时间左移位数
	sequenceMask   int64      // 序列号有效位
	nodeID         int64      // 节点ID
	lastTimestamp  int64      // 上一次生成ID的时间戳
	sequence       int64      // 序列号
}

// NewSnowflake 函数用于创建一个Snowflake实例
func NewSnowflake(nodeID int64) *Snowflake {
	s := new(Snowflake)
	s.epoch = 1729321092695 // 初始毫秒  这里是 2024-10-19
	s.nodeIDBits = 10
	s.sequenceBits = 12
	s.nodeIDShift = s.sequenceBits
	s.timestampShift = s.nodeIDBits + s.nodeIDShift
	s.sequenceMask = -1 ^ (-1 << s.sequenceBits)
	s.nodeID = nodeID
	s.lastTimestamp = -1
	s.sequence = 0
	s.mutex = sync.Mutex{}
	return s
}

// GenerateID 方法用于生成一个全局唯一的ID
func (s *Snowflake) GenerateID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentTimestamp := time.Now().UnixMilli()
	if s.lastTimestamp == currentTimestamp {
		s.sequence = (s.sequence + 1) & s.sequenceMask
		if s.sequence == 0 {
			count := 0
			for currentTimestamp <= s.lastTimestamp {
				currentTimestamp = time.Now().UnixMilli()
				time.Sleep(time.Millisecond * 2) //  时间错了休眠2毫秒
				count++
				if count > 100 {
					panic("system time is error.")
				}
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = currentTimestamp

	id := (currentTimestamp-s.epoch)<<s.timestampShift | (s.nodeID << s.nodeIDShift) | s.sequence
	return id
}
