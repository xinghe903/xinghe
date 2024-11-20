package id

import (
	"github.com/sony/sonyflake"
)

type Sonyflake struct {
	sf *sonyflake.Sonyflake
}

func NewSonyflake() *Sonyflake {
	var st sonyflake.Settings
	return &Sonyflake{sf: sonyflake.NewSonyflake(st)}
}

func (s *Sonyflake) GenerateID() uint64 {
	id, err := s.sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}
