package refreshToken

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewRefreshToken(t string, time int64) *RefreshToken {
	return &RefreshToken{
		Token:      t,
		Expiration: time,
	}
}

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(database *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: database}
}

func (r *RefreshTokenRepository) AddRefreshTokens(login string, token *RefreshToken) error {
	var err error
	_, err = r.db.Exec(`INSERT into "refreshToken" (userid, expiration, token) values ($1, $2, $3) `, login, token.Expiration, token.Token)
	return err
}

func (r *RefreshTokenRepository) GetRefreshToken(login string) (string, error) {
	refToken := RefreshToken{}
	err := r.db.QueryRowx(`SELECT userid, expiration, token FROM "refreshToken" WHERE userid = $1`, login).StructScan(&refToken)
	if err != nil {
		return "", err
	}
	return refToken.Token, nil
}
