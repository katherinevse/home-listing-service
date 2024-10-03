package middleware

import "app/pkg/auth"

type TokenManager interface {
	ParseJWT(token string) (*auth.Claims, error)
}
