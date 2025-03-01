package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Token struct {
	ID             int       `db:"id" json:"-"`
	UserID         int       `db:"user_id" json:"-"`
	TokenHash      []byte    `db:"token_hash" json:"-"`
	PlaintextToken string    `db:"plaintext_token" json:"-"`
	Expiry         time.Time `db:"expiry" json:"expiry"`
	CreatedAt      time.Time `db:"created_at" json:"-"`
}

func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
