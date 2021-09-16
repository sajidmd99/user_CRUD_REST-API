package data

import (
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Age struct {
	Value    int    `json:"value,omitempty" bson:"value,omitempty" validate:"numeric"`
	Interval string `json:"interval,omitempty" bson:"interval,omitempty" validate:"oneof=years months days weeks"`
}

type User struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Created   time.Time `json:"created,omitempty" bson:"created,omitempty"`
	Updated   time.Time `json:"updated,omitempty" bson:"updated,omitempty"`
	FirstName string    `json:"firstName,omitempty" bson:"firstName,omitempty" validate:"alpha"`
	LastName  string    `json:"lastName,omitempty" bson:"lastName,omitempty" validate:"alpha"`
	Age       Age       `json:"age,omitempty" bson:"age,omitempty" validate:"required"`
	Mobile    string    `json:"mobile,omitempty" bson:"mobile,omitempty" validate:"required,mbl"`
	Active    bool      `json:"-"`
}

func (u *User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("nm", validateName)
	validate.RegisterValidation("mbl", validateMobile)
	return validate.Struct(u)
}

func validateName(fl validator.FieldLevel) bool {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	name := fl.Field().String()
	if !isAlpha(name) {
		return false
	}

	return true
}

func validateMobile(fl validator.FieldLevel) bool {
	isNumber := regexp.MustCompile(`^[0-9]+$`).MatchString
	number := fl.Field().String()
	if !isNumber(number) || len(number) != 10 {
		return false
	}

	return true
}
