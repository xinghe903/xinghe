package token

import (
	"auth/internal/biz/auth"
	"context"
	"errors"

	"auth/internal/biz/repo"

	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
)

const (
	tokenLen = 50
)

type Token struct {
	snow     *hashid.Sonyflake
	authRepo repo.AuthRepo
	randByte *hash.RandomBytes
}

func NewToken(snow *hashid.Sonyflake, a repo.AuthRepo) auth.Auth {
	return &Token{
		snow:     snow,
		authRepo: a,
		randByte: hash.NewRandomBytes(),
	}
}

func (t *Token) GenerateToken(id string) (string, error) {
	bId, err := t.randByte.Generate(tokenLen)
	if err != nil {
		return "", errors.Join(errors.New("token generate."), err)
	}
	return string(bId), nil
}
func (t *Token) ParseToken(token string) (*auth.AuthClaims, error) {
	if len(token) == 0 {
		return nil, errors.New("token is required")
	}
	ains, err := t.authRepo.GetByToken(context.Background(), token)
	if err != nil {
		return nil, errors.Join(errors.New("token parse."), err)
	}
	return &auth.AuthClaims{
		Id:        ains.Code,
		Username:  ains.Name,
		Nickname:  ains.NickName,
		ExpiredAt: ains.ExpiredAt,
	}, nil
}
