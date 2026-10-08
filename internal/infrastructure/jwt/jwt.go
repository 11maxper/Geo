package jwt

import (
	"project/config"

	"github.com/go-chi/jwtauth/v5"
)

type Jwt struct {
	Token *jwtauth.JWTAuth
}

func NewJwt(c config.JWT) *Jwt {
	return &Jwt{
		Token: jwtauth.New("HS256", []byte(c.JWTKey), nil),
	}
}

func (j *Jwt) GetTokenString(name string) string {
	_, tokenString, _ := j.Token.Encode(map[string]interface{}{"username": name})
	return tokenString
}
