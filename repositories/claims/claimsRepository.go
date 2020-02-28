package claims

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type RepositoryOfClaims struct {
	db *sqlx.DB
}

func NewRepositoryOfClaims(database *sqlx.DB) *RepositoryOfClaims {
	return &RepositoryOfClaims{db: database}
}

func (r *RepositoryOfClaims) Create(key, value string, ctx context.Context) error {
	_, err := r.db.QueryContext(ctx, `INSERT into "claim" (key, Value) values ($1, $2)`, key, value)
	return err
}

func (r *RepositoryOfClaims) GetClaims(ctx context.Context, login string) (map[string]string, error) {
	var err error
	result := make(map[string]string, 0)
	rows, err := r.db.QueryxContext(ctx, `SELECT (key, value) FROM "claim" WHERE userid = $1`, login)
	if err != nil {
		return nil, err
	}
	claim := Claim{}
	for rows.Next() {
		err := rows.StructScan(&claim)
		//return &user, err
		_ = err

		result[claim.Description] = claim.Value
	}
	return result, err
}

func (r *RepositoryOfClaims) Delete(key, login string, ctx context.Context) error {
	_, err := r.db.QueryContext(ctx, `delete from "claim" where userid = $1 and key = $2`, login, key)
	return err
}

func (r *RepositoryOfClaims) IfExistClaim(key, login string, ctx context.Context) (bool, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT (key, value) FROM "claim" WHERE userid = $1 and key = $2 `, login, key)
	if err != nil {
		return false, nil
	}
	for rows.Next() {
		return true, err
	}
	return false, nil
}

func (r *RepositoryOfClaims) AddClaims(claims map[string]string, ctx context.Context, login string) error {
	var err error
	for k, v := range claims {
		_, err = r.db.QueryContext(ctx, `INSERT into "claim" (key, value, userid) values ($1, $2, $3) `, k, v, login)
	}
	return err
}

func (r *RepositoryOfClaims) DeleteClaims(ctx context.Context, claims map[string]string, login string) error {
	var err error
	for k, v := range claims {
		_, err = r.db.QueryContext(ctx, `DELETE FROM "claim" WHERE userid = $1 and key = $2 and value = $3`, login, k, v)
	}
	return err
}
