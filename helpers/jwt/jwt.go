package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"rakamin-final-task/helpers/errors"
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

	if claims["exp"].(int64) < time.Now().Unix() {
		return nil, errors.Unauthorized("Token expired")
	}

	if !ok {
		return nil, errors.InternalServerError("Failed to decode token")
	}

	return claims, nil
}
