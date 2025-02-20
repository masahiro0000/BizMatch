package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	RegisterInfo(user *User) error
}