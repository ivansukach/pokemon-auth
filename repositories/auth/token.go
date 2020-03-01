package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ivansukach/pokemon-auth/config"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func CreateTokenAuth(claims jwt.MapClaims) (string, error) {
	cfg := config.Load()
	claims["exp"] = time.Now().Unix() + cfg.AuthTTL
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(cfg.SecretKeyAuth))
	if err != nil {
		log.Error(err)
		return "", err
	}
	return t, nil
}

func CreateTokenRefresh(uuid string) (string, error) {
	cfg := config.Load()
	var err error
	rtClaims := jwt.MapClaims{}
	rtClaims["sub"] = "1"
	rtClaims["uuid"] = uuid
	rtClaims["exp"] = time.Now().Unix() + cfg.AuthTTL
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt, err := refreshToken.SignedString([]byte(cfg.SecretKeyRefresh))
	if err != nil {
		log.Error(err)
		return "", err
	}
	return rt, nil
}

func DecryptToken(tokS string, secretkey []byte) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokS, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return secretkey, nil
	})
	return token, err
}

func GetExpirationTimeToRefreshToken(token string, key []byte) (time.Time, error) {
	rt, err := DecryptToken(token, key)
	if err != nil {
		return time.Unix(0, 0), err
	}
	rtClaims := InterfaceToString(rt.Claims.(jwt.MapClaims))
	t, err := strconv.ParseInt(rtClaims["exp"], 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	tm := time.Unix(t, 0)
	return tm, err
}
func InterfaceToString(claims map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range claims {
		result[k] = v.(string)
	}
	return result
}
