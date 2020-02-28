package claims

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type RepositoryOfClaims struct {
	db *sqlx.DB
}

func NewRepositoryOfClaims(database *sqlx.DB) *RepositoryOfClaims {
	return &RepositoryOfClaims{db: database}
}

func (r *RepositoryOfClaims) Create(key, value string) error {
	_, err := r.db.Exec("INSERT into claim (key, value) values ($1, $2)", key, value)
	return err
}

func (r *RepositoryOfClaims) GetClaims(login string) (map[string]string, error) {
	var err error
	claims := make(map[string]string, 0)
	rows, err := r.db.Queryx("SELECT (key, value) FROM claim WHERE userid = $1", login)
	if err != nil {
		log.Warning(err)
		return nil, err
	}
	claim := Claim{}
	for rows.Next() {
		err = rows.StructScan(&claim)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		claims[claim.Description] = claim.Value
	}
	return claims, err
}

func (r *RepositoryOfClaims) Delete(key, login string) error {
	_, err := r.db.Exec("DELETE FROM claim WHERE userid = $1 AND KEY = $2", login, key)
	return err
}

func (r *RepositoryOfClaims) IfExistClaim(key, login string) (bool, error) {
	rows, err := r.db.Queryx("SELECT (key, value) FROM claim WHERE userid = $1 and key = $2 ", login, key)
	if err != nil {
		return false, nil
	}
	for rows.Next() { //WTF
		return true, err
	}
	return false, nil
}

func (r *RepositoryOfClaims) AddClaims(claims map[string]string, login string) error {
	var err error
	for k, v := range claims {
		_, err = r.db.Exec("INSERT INTO claim (key, value, userid) values ($1, $2, $3) ", k, v, login)
	}
	return err
}

func (r *RepositoryOfClaims) DeleteClaims(claims map[string]string, login string) error {
	var err error
	for k, v := range claims {
		_, err = r.db.Exec("DELETE FROM claim WHERE userid = $1 and key = $2 and value = $3", login, k, v)
	}
	return err
}
