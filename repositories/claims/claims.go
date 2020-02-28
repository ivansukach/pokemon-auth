package claims

type Claim struct {
	Description string `db:"key"`
	Value       string `db:"value"`
}
type Repository interface {
	Create(key, value string) error
	GetClaims(login string) (map[string]string, error)
	Delete(key, login string) error
	IfExistClaim(key, login string) (bool, error)
	AddClaims(claims map[string]string, login string) error
	DeleteClaims(claims map[string]string, login string) error
}
