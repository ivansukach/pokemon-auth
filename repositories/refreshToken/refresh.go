package refreshToken

type RefreshToken struct {
	UserId     string `db:"userid"`
	Token      string `db:"token"`
	Expiration int64  `db:"expiration"`
}
type Repository interface {
	AddRefreshTokens(login string, token *RefreshToken) error
	GetRefreshToken(login string) (string, error)
}
