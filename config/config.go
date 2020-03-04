package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Port             int    `env:"PORT" envDefault:"1325"`
	AuthGRPCEndpoint string `env:"AuthGRPCEndpoint" envDefault:"localhost:50051"`
	SecretKeyAuth    string `env:"SecretKeyAuth"`
	SecretKeyRefresh string `env:"SecretKeyRefresh"`
	AuthTTL          int64  `env:"TTL"`
	RefreshTTL       int64  `env:"TTL"`
}

func Load() (cfg Config) {
	cfg.AuthTTL = time.Now().Add(time.Minute*10).Unix() - time.Now().Unix()
	cfg.RefreshTTL = time.Now().Add(time.Minute*30).Unix() - time.Now().Unix()
	cfg.SecretKeyAuth = "blabla"
	cfg.SecretKeyRefresh = "blablabla"
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}
	return
}
