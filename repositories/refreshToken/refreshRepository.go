package refreshToken

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewRefreshToken(t string, time *time.Time) *RefreshToken {
	return &RefreshToken{
		Token:      t,
		Expiration: *time,
	}
}

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(database *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: database}
}

func (r *RefreshTokenRepository) AddRefreshTokens(ctx context.Context, login string, token *RefreshToken) error {
	var err error
	_, err = r.db.QueryContext(ctx, `INSERT into "refreshToken" (userid, expiration, token) values ($1, $2, $3) `, login, token.Expiration, token.Token)
	return err
}

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, login string) (string, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT userid, expiration, token FROM "refreshToken" WHERE userid = $1`, login)
	refToken := RefreshToken{}
	if err != nil {
		return "", err
	}
	for rows.Next() {
		err := rows.StructScan(&refToken)
		_ = err
	}
	return refToken.Token, err
}
