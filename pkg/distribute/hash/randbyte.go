package hash

import (
	"crypto/rand"
	"errors"
	"math"
	"math/big"
)

const (
	// 63位字符
	defaultCodes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
	hex63        = 0x3f // 63
	codeBitLen   = 6    // 63以内的数只需要6bit就能表示
	uint64BitLen = 64
)

var (
	ErrCodeLength = errors.New("codes length is not 63")
)

type RandomBytes struct {
	codes string
}

func NewRandomBytes() *RandomBytes {
	return &RandomBytes{codes: defaultCodes}
}

func (r *RandomBytes) SetCodes(codes string) error {
	if len(codes) != len(r.codes) {
		return ErrCodeLength
	}
	r.codes = codes
	return nil
}

// 随机生成指定长度字符串
func (r *RandomBytes) Generate(n int) ([]byte, error) {
	var (
		shift      = 0
		id         []byte
		randUint64 uint64
	)
	for len(id) < n {
		result, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			return nil, err
		}
		randUint64 = result.Uint64()
		shift = 0
		for shift < uint64BitLen {
			id = append(id, r.codes[randUint64&(hex63-1)])
			randUint64 >>= codeBitLen
			shift += codeBitLen
		}
	}
	return id[:n], nil
}
