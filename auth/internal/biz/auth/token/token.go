package token

import (
	"auth/internal/biz/auth"
	"context"
	"errors"
	"strconv"

	"auth/internal/biz/repo"

	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
)

const (
	tokenLen = 32
)

type Token struct {
	snow     *hashid.Snowflake
	authRepo repo.AuthRepo
}

func NewToken(snow *hashid.Snowflake, a repo.AuthRepo) auth.Auth {
	return &Token{
		snow:     snow,
		authRepo: a,
	}
}

func (t *Token) GenerateToken(id string) (string, error) {
	if len(id) == 0 {
		return "", errors.New("user id is required")
	}
	code, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", errors.Join(errors.New("token generate."), err)
	}
	snow := t.snow.GenerateID()
	token := hash.Base64Encode([]int32{int32(code >> 32), int32(code), int32(snow >> 32), int32(snow)})
	return token, nil
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
