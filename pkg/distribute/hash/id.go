package hash

import (
	hashids "github.com/speps/go-hashids/v2"
)

func GetHashId(uid int64, prefix string) string {
	hd := hashids.NewData()
	hd.Salt = "jh95xpiwren"
	hd.MinLength = 6
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err.Error())
	}
	e, err := h.Encode([]int{int(uid >> 32), int(uid & 0xFFFF)})
	if err != nil {
		panic(err.Error())
	}
	return prefix + e
}
