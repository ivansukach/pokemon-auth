package refreshToken

import (
	"context"
	"time"
)

type RefreshToken struct {
	UserId     string    `db:"userid"`
	Token      string    `db:"token"`
	Expiration time.Time `db:"expiration"`
}
type Repository interface {
	AddRefreshTokens(ctx context.Context, login string, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, login string) (string, error)
}
