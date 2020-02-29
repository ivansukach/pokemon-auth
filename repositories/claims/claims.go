package claims

import "github.com/dgrijalva/jwt-go"

type Claim struct {
	Description string `db:"key"`
	Value       string `db:"value"`
}
type Repository interface {
	GetClaims(login string) (jwt.MapClaims, error)
	IfExistClaim(key, login string) (bool, error)
	AddClaims(claims map[string]string, login string) error
	DeleteClaims(claims map[string]string, login string) error
}
