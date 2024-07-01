package jwt

import (
	"time"

	"rakamin-final-task/helpers/errors"
	"github.com/golang-jwt/jwt/v4"
)

type jwtLib struct {
	expSec    int64
	secretKey string
}

type Interface interface {
	GenerateToken(data interface{}) (string, error)
	DecodeToken(token string) (map[string]interface{}, error)
}

func Init(expSec int64, secretKey string) Interface {
	return &jwtLib{
		expSec:    expSec,
		secretKey: secretKey,
	}
}

func (j *jwtLib) GenerateToken(data interface{}) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(time.Second * time.Duration(j.expSec)).Unix(),
	})

	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtLib) DecodeToken(token string) (map[string]interface{}, error) {
	decoded, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, errors.Unauthorized("Invalid token")
	}

	if !decoded.Valid {
		return nil, errors.Unauthorized("Invalid token")
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.Unauthorized("Invalid token")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.Unauthorized("Token has expired")
	}

	return claims, nil
}
