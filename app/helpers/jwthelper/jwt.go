package jwthelper

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

var (
	//jwtSignatureKey is a secret key to hash the JWT Token
	jwtSignatureKey string
)

//AccessJWTClaims - Represent object of claims. Encourage all claims is referred to this struct
type AccessJWTClaims struct {
	jwt.StandardClaims
	Id        int   `json:"jti,omitempty"`
	ExpiresAt int64 `json:"exp,omitempty"`
}

//AccessTokenClaims - claims object for OAuth request
type AccessTokenClaims struct {
	jwt.StandardClaims
	Id         int       `json:"id"`
	Scopes     []string  `json:"scopes"`
	ExpiresAt  int64     `json:"exp,omitempty"`
	ClientUUID uuid.UUID `json:"client_uuid,omitempty"`
}

func init() {
	err := godotenv.Load("application.env")
	if err != nil {
		log.Error().Msg("Failed read configuration database")
		return
	}
	jwtSignatureKey = os.Getenv("jwt.secretKey")
}

//NewWithClaims will return token with custom claims
func NewWithClaims(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSignatureKey))

	if err != nil {
		return "", err
	}
	return ss, nil
}

//VerifyTokenWithClaims will verify the validity of token and return the claims
func VerifyTokenWithClaims(token string) (*AccessJWTClaims, error) {

	jwtToken, err := jwt.ParseWithClaims(token, &AccessJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSignatureKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*AccessJWTClaims)
	if !ok {
		return nil, errors.New("error retrieving claims")
	}

	timeNow := time.Now().Unix()
	if claims.ExpiresAt < timeNow {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

func VerifyGrantChallengeToken(token string) (*AccessTokenClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &AccessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSignatureKey), nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, errors.New("error retrieving claims")
	}

	timeNow := time.Now().Unix()
	if claims.ExpiresAt < timeNow {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

func generateAccessToken(userId int, scopes []string, durationMinute int) (token string, err error) {
	// Set expiration time to 10 minutes
	minute := 60 //in second
	jwtExpiredAt := time.Now().Unix() + int64(10*minute)

	challengeClaims := AccessTokenClaims{
		Id:        userId,
		Scopes:    scopes,
		ExpiresAt: jwtExpiredAt,
	}

	jwtToken, err := NewWithClaims(challengeClaims)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
