package tools

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type JwtClaim struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func TokenCreate(id int) string {
	var jwtKey = []byte(os.Getenv("JWT_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaim{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	signedStr, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return signedStr
}

func TokenValidate(t string) (*jwt.Token, error) {
	var jwtKey = []byte(os.Getenv("JWT_KEY"))
	token, _ := jwt.ParseWithClaims(t, &JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			panic(gqlerror.Errorf("Error decoding token"))
		}

		return jwtKey, nil
	})

	return token, nil
}
