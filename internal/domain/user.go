package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

// User represents a user account in the system.
type User struct {
	ID 					int64  `db:"id"`
	Username 			string `db:"username"`
	DisplayName			string `db:"display_name"`
	Password 			string `db:"password"`
	Prefecture			*int64 `db:"prefecture_id"`
	PrefectureName		*string `db:"prefecture"`
	Industry			*int64 `db:"industry_id"`
	IndustryName		*string `db:"industry"`
	Job					*int64 `db:"job_id"`
	JobName				*string `db:"job"`
	Position 			*int64 `db:"position_id"`
	PositionName		*string `db:"position"`
	Age 				*int64 `db:"age"`
	Gender 				*string `db:"gender"`
	Photo				*string `db:"photo"`
	ProfileDescription	*string `db:"profile_description"`
}

// UserSearchFilter holds criteria used to filter and search for specific users.
type UserSearchFilter struct {
	Prefectures		[]*int64 `db:"prefecture_id"`
	Industries		[]*int64 `db:"industry_id"`
	Jobs			[]*int64 `db:"job_id"`
	Positions		[]*int64 `db:"position_id"`
	Ages			[]*int64 `db:"age"`
	AgeGroups		[]int
	Genders			[]*string `db:"gender"`
	ExcludeUserID	*int64
}

// ScoredUser represents a user with a calculated score based on the current user's information.
type ScoredUser struct {
	User *User
	Score int
}

var (
	ErrInvalidInput		 	= errors.New("入力内容が不正です")
	ErrPasswordTooShort  	= errors.New("パスワードは8文字以上必要です")
	ErrUserAlreadyExists 	= errors.New("ユーザーは既に存在します")
	ErrIncorrectPassword 	= errors.New("パスワードが一致しません")
	ErrOldPasswordMismatch	= errors.New("古いパスワードが一致しません")
	ErrNewPasswordMismatch	= errors.New("新しいパスワードが一致しません")
	ErrUserNotFound			= errors.New("ユーザーが見つかりません")
)

// Create a new user instance by validating the inputs and hashing the password.
func NewUser(username, displayName, password string) (*User, error) {
	// Validate that the username and display name are not empty.
	if username == "" || displayName == "" {
		return nil, ErrInvalidInput
	}

	// Check if the password meets the minimum length requirement.
	if len(password) < MinPasswordLength {
		return nil, ErrPasswordTooShort
	}

	// Generate a bcrypt hash of the password.
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Username: 		username,
		DisplayName: 	displayName,
		Password:		hashedPassword,
	}, nil
}

// hashPassword takes a plaintext password and returns its bcrypt hash.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword compares the provided password with the stored hash.
func (u *User) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return ErrIncorrectPassword
	}
	return nil
}

// CheckUsernameDisplayName validates that both username and display name are not empty.
func CheckUsernameDisplayName(username string, displayName string) error {
	if username == "" || displayName == "" {
		return ErrInvalidInput
	}
	return nil
}