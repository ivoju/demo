package jwt

import (
	"errors"
	"strconv"
	"strings"
	"time"

	actpb "github.com/demo/proto/v1.0/accounts"
	"github.com/dgrijalva/jwt-go"
	"github.com/kenshaw/envcfg"
	"google.golang.org/grpc"
)

const (
	ISSUER = `Bank Raya Indonesia`
	TYPE   = `Bearer`
)

type CustomClaims struct {
	UserId string `json:"userId,omitempty"`
	jwt.StandardClaims
}

type Token struct {
	Type       string
	Access     string
	ExpPeriode uint32
}

type Config struct {
	Env *envcfg.Envcfg
}

// func to verify fixed token
func IsFixedTokenVerified(config *envcfg.Envcfg, auth string) error {
	// Bearer token as  RFC 6750 standard
	if strings.Split(auth, " ")[0] != TYPE || strings.Split(auth, " ")[1] != config.GetKey("jwt.fixedtoken") {
		return errors.New("Invalid token")
	}

	return nil
}

func New(env *envcfg.Envcfg) *Config {
	return &Config{
		Env: env,
	}
}

func (conf *Config) GenerateToken(userid string) (*actpb.Token, error) {
	expPeriode, err := strconv.Atoi(conf.Env.GetKey("jwt.expperiode"))
	if err != nil {
		return nil, err
	}

	expPeriodeDuration := time.Duration(expPeriode)

	// Create the Claims
	claims := CustomClaims{
		UserId: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expPeriodeDuration * time.Second).Unix(),
			Issuer:    ISSUER,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(conf.Env.GetKey("jwt.key")))
	if err != nil {
		return nil, err
	}

	return &actpb.Token{
		Type:       TYPE,
		Access:     accessToken,
		ExpPeriode: uint32(expPeriodeDuration),
	}, nil
}

func (conf *Config) ClaimToken(auth string, info *grpc.UnaryServerInfo) (jwt.MapClaims, error) {
	// Bearer token as  RFC 6750 standard
	if strings.Split(auth, " ")[0] != TYPE {
		return nil, errors.New("Invalid token")
	}

	token, err := jwt.Parse(strings.Split(auth, " ")[1], func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errors.New("Signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			errors.New("Signing method invalid")
		}

		return []byte(conf.Env.GetKey("jwt.key")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("failed to claim token")
	}

	return claims, nil
}
