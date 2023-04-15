package utils

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTsecret = []byte("BTBT")

type Claims struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username string, uid uint, password string) (string, error) {
	nowtTime := time.Now()
	expireTime := nowtTime.Add(24 * time.Hour)
	myClaims := Claims{
		UserName: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  nowtTime.Unix(),
			ExpiresAt: expireTime.Unix(),
			Id:        strconv.Itoa(int(uid)),
		},
	}
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	token, err := tokenClaim.SignedString(JWTsecret)
	return token, err
}
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTsecret, nil
	})
	if tokenClaims != nil {
		if myClaims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return myClaims, nil
		}
	}
	return nil, err

}
