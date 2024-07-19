package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	IssuerName      string `json:"IssuerName"`
	JwtSignatureKey []byte `json:"JwtSignatureKey"`
}

type TokenConfig struct {
	SecretKey string
	Issuer    string
}

func NewJwtService(config *TokenConfig) *JwtService {
	return &JwtService{
		IssuerName:      config.Issuer,
		JwtSignatureKey: []byte(config.SecretKey),
	}
}

func (s *JwtService) GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
		"iss":    s.IssuerName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JwtSignatureKey)
}

func (s *JwtService) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.JwtSignatureKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
