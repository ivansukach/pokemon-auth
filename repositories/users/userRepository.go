package users

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type userRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(user *User) error {
	_, err := ur.db.NamedExec("INSERT INTO users VALUES (:login, :password, :name, :surname, :coins)", user)
	return err
}
func (ur *userRepository) Get(login string) (*User, error) {
	u := User{}
	err := ur.db.QueryRowx("SELECT * FROM users WHERE Login=$1", login).StructScan(&u)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &u, err
}
func (ur *userRepository) Update(user *User) error {
	_, err := ur.db.NamedExec("UPDATE users SET (Login=:login, Password=:password"+
		"Name=:name, Surname=:surname, Coins=:coins) WHERE Login=:login", user)
	return err
}
func (ur *userRepository) Delete(login string) error {
	_, err := ur.db.Exec("DELETE FROM users WHERE login=$1", login)
	return err
}
func (ur *userRepository) Listing() ([]User, error) {
	rows, err := ur.db.Queryx("SELECT * FROM users")
	if err != nil {
		log.Warning(err)
		return nil, err
	}
	u := make([]User, 0)
	for rows.Next() {
		user := User{}
		err = rows.StructScan(&user)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		u = append(u, user)
	}
	return u, err
}
