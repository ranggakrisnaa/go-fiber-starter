package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateAccessToken(userId string, roles []string) string
	GenerateRefreshToken() (string, time.Time)
	ValidateToken(token string) (*jwt.Token, error)
	GetUserIDByToken(token string) (string, error)
	GetRolesByToken(token string) ([]string, error)
}

type jwtCustomClaim struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey     string
	issuer        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey:     getSecretKey(),
		issuer:        "Template",
		accessExpiry:  time.Minute * 15,
		refreshExpiry: time.Hour * 24 * 7,
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "Template"
	}
	return secretKey
}

func (j *jwtService) GenerateAccessToken(userId string, roles []string) string {
	claims := jwtCustomClaim{
		userId,
		roles,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tx, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err)
	}
	return tx
}

func (j *jwtService) GenerateRefreshToken() (string, time.Time) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return "", time.Time{}
	}

	refreshToken := base64.StdEncoding.EncodeToString(b)
	expiresAt := time.Now().Add(j.refreshExpiry)

	return refreshToken, expiresAt
}

func (j *jwtService) parseToken(t_ *jwt.Token) (any, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
	}
	return []byte(j.secretKey), nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, j.parseToken)
}

func (j *jwtService) GetUserIDByToken(token string) (string, error) {
	tToken, err := j.ValidateToken(token)
	if err != nil {
		return "", err
	}

	claims := tToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id, nil
}

func (j *jwtService) GetRolesByToken(token string) ([]string, error) {
	tToken, err := j.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	claims := tToken.Claims.(jwt.MapClaims)
	rolesInterface, ok := claims["roles"]
	if !ok {
		return nil, fmt.Errorf("roles not found in token")
	}

	rolesSlice, ok := rolesInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("roles format is invalid")
	}

	roles := make([]string, len(rolesSlice))
	for i, role := range rolesSlice {
		roleStr, ok := role.(string)
		if !ok {
			return nil, fmt.Errorf("role at index %d is not a string", i)
		}
		roles[i] = roleStr
	}

	return roles, nil
}
