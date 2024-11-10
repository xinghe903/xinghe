package auth

import "time"

type AuthClaims struct {
	Id        string
	Username  string
	Nickname  string
	ExpiredAt time.Time
}

type Auth interface {
	GenerateToken(id string) (string, error)
	ParseToken(token string) (*AuthClaims, error)
}
