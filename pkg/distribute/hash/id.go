package hash

import (
	"encoding/base32"
	"encoding/base64"
	"strings"

	hashids "github.com/speps/go-hashids/v2"
)

func GetHashId(uid int64, prefix string) string {
	e := GetHashIdWithIds([]int{int(uid >> 32), int(uid & 0xFFFF)}, 6)
	return prefix + e
}

func GetHashIdWithIds(ids []int, mLen int) string {
	hd := hashids.NewData()
	hd.Salt = "jh95xpiwren"
	hd.MinLength = mLen
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err.Error())
	}
	e, err := h.Encode(ids)
	if err != nil {
		panic(err.Error())
	}
	return e
}

func Base64Encode(targets []int32) string {
	var bases []byte
	for _, target := range targets {
		bases = append(bases, byte(target>>24), byte(target>>16), byte(target>>8), byte(target))
	}
	encode := base64.StdEncoding.EncodeToString(bases)
	return strings.TrimRight(encode, "=")
}

func Base32Encode(targets []int32) string {
	var bases []byte
	for _, target := range targets {
		bases = append(bases, byte(target>>24), byte(target>>16), byte(target>>8), byte(target))
	}
	encode := base32.HexEncoding.EncodeToString(bases)
	return strings.TrimRight(encode, "=")
}
