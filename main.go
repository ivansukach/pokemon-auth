package main

import (
	"github.com/ivansukach/pokemon-auth/protocol"
	"github.com/ivansukach/pokemon-auth/repositories/claims"
	"github.com/ivansukach/pokemon-auth/repositories/refreshToken"
	"github.com/ivansukach/pokemon-auth/repositories/users"
	"github.com/ivansukach/pokemon-auth/server"
	"github.com/ivansukach/pokemon-auth/service"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=su password=su dbname=pokemons sslmode=disable")
	rpsClaims := claims.NewRepositoryOfClaims(db)
	rpsUsers := users.New(db)
	rpsRefreshTokens := refreshToken.NewRefreshTokenRepository(db)
	log.Info("GRPC-server started")
	s := grpc.NewServer()
	as := service.New(rpsUsers, rpsClaims, rpsRefreshTokens)

	srv := server.New(as)
	protocol.RegisterAuthServiceServer(s, srv)  //зарегистрировали сервер
	listener, err := net.Listen("tcp", ":1325") //просто слушаем порт
	s.Serve(listener)
	if err != nil {
		log.Error(err)
	}

	defer db.Close()
}
