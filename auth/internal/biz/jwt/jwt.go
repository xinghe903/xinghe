package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	JWTSigningKey = "JWT_SIGNING_KEY"
	JWTIssuer     = "JWT_ISSUER"
	JWTBufferTime = 1800 // seconds

	JWTSigningKeyDBValue = "oFMW0t3pL4ezfzq8"
)

var (
	ErrTokenParse = errors.New("ErrTokenParse")
)

type JWT struct {
	SigningKey []byte
}

type JwtRequest struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	ID   uint
	Name string
}

func NewJWT() *JWT {
	// settingRepo := repo.NewISettingRepo()
	// jwtSign, _ := settingRepo.Get(settingRepo.WithByKey("JWTSigningKey"))
	// return &JWT{
	// 	[]byte(jwtSign.Value),
	// }
	return &JWT{
		[]byte(JWTSigningKeyDBValue),
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: JWTBufferTime,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(JWTBufferTime))),
			Issuer:    JWTIssuer,
		},
	}
	return claims
}

func (j *JWT) CreateToken(request CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &request)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenStr string) (*JwtRequest, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtRequest{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil || token == nil {
		return nil, ErrTokenParse
	}
	if claims, ok := token.Claims.(*JwtRequest); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenParse
}

func GenerateToken(name string) string {
	j := NewJWT()
	claims := j.CreateClaims(BaseClaims{
		Name: name,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated Token:", token)
	return token
}

func ParseToken(tokenStr string) *CustomClaims {
	j := NewJWT()
	claims, err := j.ParseToken(tokenStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed Claims:", claims)
	// fmt.Println(claims.ExpiresAt.)
	return (*CustomClaims)(claims)
}
