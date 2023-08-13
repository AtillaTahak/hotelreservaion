package types

import (
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minPasswordLength  = 8
	minFirstNameLength = 2
	minLastNameLength  = 2
	maxFirstNameLength = 50
	maxLastNameLength  = 50
	maxEmailLength     = 254
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (params *UpdateUserParams) ToBSON() bson.M {
	update := bson.M{}
	if params.FirstName != "" {
		update["firstName"] = params.FirstName
	}
	if params.LastName != "" {
		update["lastName"] = params.LastName
	}
	return update
}
func (params *CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Password) < minPasswordLength {
		errors["Password"] = "password is too short"
	}
	if len(params.FirstName) < minFirstNameLength {
		errors["firstName"] = "first name is too short"
	}
	if len(params.LastName) < minLastNameLength {
		errors["lastName"] = "last name is too short"
	}
	if len(params.FirstName) > maxFirstNameLength {
		errors["firstName"] = "first name is too long"
	}
	if len(params.LastName) > maxLastNameLength {
		errors["lastName"] = "last name is too long"
	}
	if len(params.Email) > maxEmailLength {
		errors["email"] = "email is too long"
	}
	if !isEmailValid(params.Email) {
		errors["email"] = "email is invalid"
	}
	return errors
}
func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}
func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9` +
		`-]+(?:\.[a-zA-Z0-9` + `-]+)*$`)
	return emailRegex.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

func NewUserFromParams(params *CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
