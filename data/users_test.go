package data

import "testing"

func TestCheckValidation(t *testing.T) {

	u := &User{
		FirstName: "Walter",
		LastName:  "White",
		Age: Age{
			Value:    50,
			Interval: "years",
		},
		Mobile: "9823574857",
	}

	err := u.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
