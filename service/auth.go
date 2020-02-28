package service

import (
	"github.com/ivansukach/pokemon-auth/repositories"
	"github.com/ivansukach/pokemon-auth/repositories/claims"
	"github.com/ivansukach/pokemon-auth/repositories/refreshToken"
	"github.com/ivansukach/pokemon-auth/repositories/users"
)

type UserService struct {
	ur      users.Repository
	claims  claims.RepositoryOfClaims
	refresh refreshToken.RefreshTokenRepository
}

func New(repo claims2.Repository) *UserService {
	return &UserService{r: repo}
}
