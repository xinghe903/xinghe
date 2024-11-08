package id

import (
	"fmt"
	"net"
	"strings"
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
// sgin   timestamp    nodeId   sequence
//
//	1        41          10        12
//
// 符号位1  时间戳41  节点ID10  序号12  1+41+10+12=64位
func NewSnowflake(dns string) *Snowflake {
	s := new(Snowflake)
	s.epoch = 1729321092695 // 初始毫秒  这里是 2024-10-19
	s.nodeIDBits = 10
	s.sequenceBits = 12
	s.nodeIDShift = s.sequenceBits
	s.timestampShift = s.nodeIDBits + s.nodeIDShift
	s.sequenceMask = -1 ^ (-1 << s.sequenceBits)
	s.nodeID = int64(getLocalIpv4Uint32(dns)) & int64(GetNodeIdByBitCnt(s.nodeIDBits)) // 保留 nodeIDBits 位
	fmt.Printf("nodeId: %d\n", s.nodeID)
	s.lastTimestamp = -1
	s.sequence = 0
	s.mutex = sync.Mutex{}
	return s
}

func getLocalIpv4Uint32(dns string) uint32 {
	if dns == "" {
		dns = "8.8.8.8:53"
	}
	conn, err := net.Dial("udp", dns)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	addr := conn.LocalAddr().(*net.UDPAddr)
	ipstr := strings.Split(addr.String(), ":")[0]
	var ipuint32 uint32 = uint32(addr.IP[0])<<24 | uint32(addr.IP[1])<<16 | uint32(addr.IP[2])<<8 | uint32(addr.IP[3])
	fmt.Printf("Local IP: %s, uint32: %d\n", ipstr, ipuint32)
	return ipuint32
}

func GetNodeIdByBitCnt(bits uint) uint64 {
	var bnum16 uint64 = 0
	for bits > 0 {
		bnum16 = bnum16 | 1<<(bits-1)
		bits--
	}
	return bnum16
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
