package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/maragudk/go-ahead/errors2"
	"github.com/maragudk/go-ahead/model"
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
// Note that the password comparison is done no matter what, to combat timing attacks.
// See https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#compare-password-hashes-using-safe-functions
func (s *Storer) Login(ctx context.Context, email, password string) (*model.User, error) {
	// Start by filling in a password, otherwise bcrypt.CompareHashAndPassword returns bcrypt.ErrHashTooShort when no user was found.
	var user model.User
	if err := s.DB.Get(&user, loginQuery, email); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		// If no user was found, add a random password, or bcrypt.CompareHashAndPassword below
		// returns early with bcrypt.ErrHashTooShort
		user.Password = "$2a$10$0flYAPknlgvCeBICHGFKWeMeWRa3ZDcKXVihym71oPXNDcCj//8LC"
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil
	}
	user.Password = ""
	return &user, nil
}
