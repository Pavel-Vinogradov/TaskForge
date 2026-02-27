package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int) (string, error) {
	viper.SetDefault("jwt.secret", "cutytyhbheq")
	viper.SetDefault("jwt.expiration", "24h")

	jwtSecret := viper.GetString("jwt.secret")
	expiration := viper.GetString("jwt.expiration")

	duration, err := time.ParseDuration(expiration)
	if err != nil {
		duration = 24 * time.Hour
	}

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	viper.SetDefault("jwt.secret", "cutytyhbheq")
	jwtSecret := viper.GetString("jwt.secret")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
