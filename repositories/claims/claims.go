package claims

import "context"

type Claim struct {
	Description string `db:"key"`
	Value       string `db:"value"`
}
type Repository interface {
	Create(key, value string, ctx context.Context) error
	GetClaims(ctx context.Context, login string) (map[string]string, error)
	Delete(key, login string, ctx context.Context) error
	IfExistClaim(key, login string, ctx context.Context) (bool, error)
	AddClaims(claims map[string]string, ctx context.Context, login string) error
	DeleteClaims(ctx context.Context, claims map[string]string, login string) error
}
