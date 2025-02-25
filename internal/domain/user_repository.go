package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(userID int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	RegisterInfo(user *User) error
	UpdatePassword(userID int64, hashedPassword string) error
	SearchUsers(filter *UserSearchFilter) ([]*User, error)
}