package auth

type Repository interface {
	CreateUser(email string, password []byte) error
	GetUserByEmail(email string) (*User, error)
}
