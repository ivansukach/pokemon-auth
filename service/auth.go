package service

import (
	"errors"
	"fmt"
	"github.com/ivansukach/pokemon-auth/repositories/claims"
	"github.com/ivansukach/pokemon-auth/repositories/refreshToken"
	"github.com/ivansukach/pokemon-auth/repositories/users"
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
func (gs *Auth) CreateUser(user *users.User) error {
	return gs.ur.Create(user)
}
func (gs *Auth) UpdateUser(user *users.User) error {
	return gs.ur.Update(user)
}
func (gs *Auth) GetUser(login string) (*users.User, error) {
	return gs.ur.Get(login)
}
func (gs *Auth) DeleteUser(id string) error {
	return gs.ur.Delete(id)
}
func (gs *Auth) ListingUsers() ([]users.User, error) {
	return gs.ur.Listing()
}

func (gs *Auth) DeleteClaims(claims map[string]string, login string) error {
	return gs.claims.DeleteClaims(ctx, claims, login)
}

func (gs *Auth) AddClaims(claims map[string]string, login string) error {
	return gs.claims.AddClaims(claims, ctx, login)
}
func (s *UserService) RefreshToken(ctx context.Context, tokenReqAuth, tokenReqRefresh string) (string, string, error) {
	token, err := auth.GetTokenAuth(tokenReqAuth)
	if err != nil {
		return "", "", err
	}
	var newToken, newRefToken string

	claims := repository.IntefaceToString(token.Claims.(jwt.MapClaims))

	tokenRefresh, err := auth.GetTokenRefresh(tokenReqRefresh)
	if err != nil {
		return "", "", err
	}

	claimsR := repository.IntefaceToString(tokenRefresh.Claims.(jwt.MapClaims))

	if claims["login"] == claimsR["uuid"] {

		user, err := s.users.FindUser(ctx, claims["login"])
		if err != nil {
			log.Errorf("User not found when refresh token", err)
			return "", "", err
		}
		refreshToken, err := s.refresh.GetRefreshToken(ctx, claims["login"])

		if ok := s.users.Delete(ctx, claims["login"]); ok != nil {
			return "", "", ok
		}

		tm, err := auth.GetExpirationTimeToRefreshToken(refreshToken)

		if time.Now().After(tm) {
			log.Errorf("Invalid session")
		}

		err = s.users.Create(ctx, user)
		if err != nil {
			return "", "", err
		}

		newClaims, err := s.claims.GetClaims(ctx, claims["login"])
		if err != nil {
			return "", "", nil
		}
		uuid := guuid.New().String()
		newToken, err = auth.CreatTokenAuth(claims["login"], newClaims)
		newRefToken, err = auth.CreatTokenRefresh(uuid)
		if err != nil {
			return "", "", err
		}

		t, err := auth.GetExpirationTimeToRefreshToken(newRefToken)
		if err != nil {
			return "", "", err
		}
		err = s.refresh.AddRefreshTokens(ctx, claims["login"], repository.NewRefreshToken(newRefToken, &t))
	} else {
		return "", "", errors.New("not validate in refresh")
	}

	return newToken, newRefToken, err
}
func New(usersRepository repository.UsersRepository, tokenRepository repository.RefreshTokenRepository, claims repository.RepositoryOfClaims, config conf.Config) *UserService {
	return &UserService{
		users:   &usersRepository,
		claims:  &claims,
		refresh: &tokenRepository,
		cfg:     &config,
	}
}

func (s *UserService) SignIn(ctx context.Context, login string, password string) (token string, tokenRefresh string, err error) {
	user, err := s.users.FindUser(ctx, login)

	if err != nil {
		return "", "", err
	}
	if user != nil {
		if err := repository.VerifyPassword(string(user.Password), password); user.Login == login && err == nil {

			claims, err := s.claims.GetClaims(ctx, login)
			if err != nil {
				return "", "", err
			}

			token, err = auth.CreatTokenAuth(login, claims)
			if err != nil {
				return "", "", err
			}

			tokenRefresh, err = auth.CreatTokenRefresh(login)
			if err != nil {
				return "", "", err
			}
		}
	} else {
		return "", "", err
	}

	return token, tokenRefresh, err
}

func (s *UserService) SignUp(ctx context.Context, login string, password string) (string, error) {
	var err error

	if login == "" || password == "" {
		return "", errors.New("verify user")
	}

	if s.users.IfExistUser(ctx, login) {
		return "", fmt.Errorf("enter other login")
	}

	resp, err := SendgridMsg(s.cfg.SendGridKey)

	if err != nil {
		return "", err
	}

	uuid := guuid.New().String()
	claims, err := s.claims.GetClaims(ctx, login)
	if err != nil {
		return "", nil
	}
	token, err := auth.CreatTokenAuth(login, claims)
	refreshToken, err := auth.CreatTokenRefresh(uuid)

	if err != nil {
		log.Errorf("token not created in sign up", err)
		return "", err
	}

	err = s.users.Create(ctx, &repository.User{
		Login:    login,
		Password: repository.Hash(password),
	})

	if err != nil {
		return "", err
	}

	t, err := auth.GetExpirationTimeToRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	err = s.refresh.AddRefreshTokens(ctx, login, repository.NewRefreshToken(refreshToken, &t))

	url := fmt.Sprintf("http://localhost:%d/confirm?token=%s", s.cfg.Port, s.cfg.SecretKeyAuth)

	return url, err
}
