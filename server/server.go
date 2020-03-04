package server

import (
	"context"
	"github.com/ivansukach/pokemon-auth/protocol"
	"github.com/ivansukach/pokemon-auth/repositories/users"
	"github.com/ivansukach/pokemon-auth/service"
	log "github.com/sirupsen/logrus"
)

//Server ...
type Server struct {
	as *service.Auth
}

func New(as *service.Auth) *Server {
	return &Server{as: as}
}
func (s *Server) AddClaims(ctx context.Context, req *protocol.AddClaimsRequest) (*protocol.EmptyResponse, error) {
	err := s.as.AddClaims(req.Claims, req.Login)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *protocol.DeleteRequest) (*protocol.EmptyResponse, error) {
	err := s.as.DeleteUser(req.Login)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}

func (s *Server) DeleteClaims(ctx context.Context, req *protocol.DeleteClaimsRequest) (*protocol.EmptyResponse, error) {
	err := s.as.DeleteClaims(req.Claims, req.Login)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *protocol.RefreshTokenRequest) (*protocol.RefreshTokenResponse, error) {
	token, refToken, err := s.as.RefreshToken(req.Token, req.TokenRefresh)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	response := &protocol.RefreshTokenResponse{
		Token:        token,
		RefreshToken: refToken,
	}
	return response, nil
}

func (s *Server) SignIn(ctx context.Context, req *protocol.SignInRequest) (*protocol.SignInResponse, error) {
	token, refresh, err := s.as.SignIn(req.Login, req.Password)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	user, err := s.as.GetUser(req.Login)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	response := &protocol.SignInResponse{
		Name:         user.Name,
		Surname:      user.Surname,
		Coins:        user.Coins,
		Photo:        user.Photo,
		Login:        user.Login,
		Token:        token,
		RefreshToken: refresh,
	}
	return response, nil
}

func (s *Server) SignUp(ctx context.Context, req *protocol.SignUpRequest) (*protocol.EmptyResponse, error) {
	user := users.User{Login: req.Login, Password: req.Password, Name: req.Name, Surname: req.Surname, Coins: 10000}
	err := s.as.SignUp(&user)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}
