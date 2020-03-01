package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ivansukach/pokemon-auth/config"
	"github.com/ivansukach/pokemon-auth/repositories/auth"
	"github.com/ivansukach/pokemon-auth/repositories/claims"
	"github.com/ivansukach/pokemon-auth/repositories/refreshToken"
	"github.com/ivansukach/pokemon-auth/repositories/users"
	log "github.com/sirupsen/logrus"
	"time"
)

type Auth struct {
	ur      users.Repository
	claims  claims.Repository
	refresh refreshToken.Repository
}

func New(ur users.Repository, claims claims.Repository, refresh refreshToken.Repository) *Auth {
	return &Auth{ur: ur, claims: claims, refresh: refresh}
}
func (as *Auth) CreateUser(user *users.User) error {
	return as.ur.Create(user)
}
func (as *Auth) UpdateUser(user *users.User) error {
	return as.ur.Update(user)
}
func (as *Auth) GetUser(login string) (*users.User, error) {
	return as.ur.Get(login)
}
func (as *Auth) DeleteUser(id string) error {
	return as.ur.Delete(id)
}
func (as *Auth) ListingUsers() ([]users.User, error) {
	return as.ur.Listing()
}

func (as *Auth) DeleteClaims(claims map[string]string, login string) error {
	return as.claims.DeleteClaims(claims, login)
}

func (as *Auth) AddClaims(claims map[string]string, login string) error {
	return as.claims.AddClaims(claims, login)
}

func (as *Auth) RefreshToken(tokenReqAuth, tokenReqRefresh string) (newToken string, newRefToken string, err error) {
	cfg := config.Load()
	token, err := auth.DecryptToken(tokenReqAuth, []byte(cfg.SecretKeyAuth))
	if err != nil {
		return
	}
	tokenRefresh, err := auth.DecryptToken(tokenReqRefresh, []byte(cfg.SecretKeyRefresh))
	if err != nil {
		return
	}

	claims := auth.InterfaceToString(token.Claims.(jwt.MapClaims))
	claimsR := auth.InterfaceToString(tokenRefresh.Claims.(jwt.MapClaims))
	login := claims["login"]
	if claims["login"] == claimsR["uuid"] {
		rToken, err := as.refresh.GetRefreshToken(login)
		decryptedToken, err := auth.DecryptToken(rToken, []byte(cfg.SecretKeyRefresh))
		claims := decryptedToken.Claims.(jwt.MapClaims)
		uniqueIdentifier := uuid.New().String()
		newToken, err = auth.CreateTokenAuth(claims)
		if time.Now().Unix() > claims["exp"].(int64) {
			log.Errorf("Invalid session")
			newRefToken, err = auth.CreateTokenRefresh(uniqueIdentifier)
		}
		if err != nil {
			return "", "", err
		}
		err = as.refresh.AddRefreshTokens(login, refreshToken.NewRefreshToken(newRefToken, claims["exp"].(int64)))
	} else {
		return "", "", errors.New("Fake refresh token ")
	}
	return
}
func (as *Auth) SignIn(login string, password string) (token string, tokenRefresh string, err error) {
	user, err := as.ur.Get(login)

	if err != nil {
		return
	}
	if password == user.Password {
		claims, err := as.claims.GetClaims(login)
		if err != nil {
			return "", "", err
		}
		token, err = auth.CreateTokenAuth(claims)
		if err != nil {
			return "", "", err
		}
		tokenRefresh, err = auth.CreateTokenRefresh(login)
		if err != nil {
			return "", "", err
		}
	}
	return
}

func (as *Auth) SignUp(user *users.User) error {
	if user.Login == "" || user.Password == "" {
		return errors.New("Empty fields!!! ")
	}
	_, err := as.ur.Get(user.Login)
	if user == nil || err == nil {
		return fmt.Errorf("enter other login")
	}
	err = as.ur.Create(user)
	claims := make(map[string]string, 0)
	claims["login"] = user.Login
	err = as.claims.AddClaims(claims, user.Login)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
