package users

type User struct {
	Login    string `db:"login"`
	Password string `db:"password"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
	Coins    int32  `db:"coins"`
}
type Repository interface {
	Create(user *User) error
	Get(login string) (*User, error)
	Update(user *User) error
	Delete(login string) error
	Listing() ([]User, error)
}
