package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"go-ahead/errors2"
	"go-ahead/model"
)

const (
	signupQuery = "insert into users (name, email, password) values (:name, :email, :password)"
	loginQuery  = "select id, name, password, email, created, updated from users where email = $1"
)

const (
	minPasswordLength = 10
	maxPasswordLength = 64
)

// Signup a user. Note that the given password is in cleartext and hashed here.
func (s *Storer) Signup(ctx context.Context, name, email, password string) error {
	passwordLength := utf8.RuneCountInString(password)
	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		return fmt.Errorf("%v is outside the password length range of [%v,%v]", passwordLength, minPasswordLength, maxPasswordLength)
	}

	if !model.ValidateEmail(email) {
		return fmt.Errorf("%v is not an email address", email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors2.Wrap(err, "could not hash password")
	}
	user := model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	if _, err := s.DB.NamedExecContext(ctx, signupQuery, user); err != nil {
		return err
	}
	return nil
}

// Login with a given email and password. The password is cleartext and hashed here.
// Returns the user if succesful, without the password.
func (s *Storer) Login(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User
	if err := s.DB.Get(&user, loginQuery, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil
	}
	user.Password = ""
	return &user, nil
}
