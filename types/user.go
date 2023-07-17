package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)
const (
	bcryptCost = 12
	minPasswordLength = 8
	minFirstNameLength = 2
	minLastNameLength = 2
	maxFirstNameLength = 50
	maxLastNameLength = 50
	maxEmailLength = 254
)
type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params *CreateUserParams) Validate() []error {
	errors := []error{}
	if len(params.Password) < minPasswordLength {
		errors = append(errors,fmt.Errorf("password is too short"))
	}
	if len(params.FirstName) < minFirstNameLength {
		errors = append(errors,fmt.Errorf("first name is too short"))
	}
	if len(params.LastName) < minLastNameLength {
		errors = append(errors,fmt.Errorf("last name is too short"))
	}
	if len(params.FirstName) > maxFirstNameLength {
		errors = append(errors,fmt.Errorf("first name is too long"))
	}
	if len(params.LastName) > maxLastNameLength {
		errors = append(errors,fmt.Errorf("last name is too long"))
	}
	if len(params.Email) > maxEmailLength {
		errors = append(errors,fmt.Errorf("email is too long"))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors,fmt.Errorf("email is invalid"))
	}
	return errors
}
func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9` +
		`-]+(?:\.[a-zA-Z0-9` + `-]+)*$`)
	return emailRegex.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string `bson:"firstName" json:"firstName"`
	LastName          string `bson:"lastName" json:"lastName"`
	Email             string `bson:"email" json:"email"`
	EncryptedPassword string `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params *CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil,err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}