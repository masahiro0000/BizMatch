package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

type User struct {
	Username 		string
	DisplayName 	string
	Password 		string
}

var (
	ErrInvalidInput		 = errors.New("入力内容が不正です")
	ErrPasswordTooShort  = errors.New("パスワードは8文字以上必要です")
	ErrUserAlreadyExists = errors.New("ユーザーは既に存在します")
	ErrIncorrectPassword = errors.New("パスワードが一致しません")
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
	hashedPassword, err := hashPassword(password)
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
func hashPassword(password string) (string, error) {
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