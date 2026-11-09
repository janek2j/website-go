package models

import (
	"errors"
	"time"
)

type Person struct {
	Id        string
	Name      string
	Email     string
	Phone     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Person) String() string {
	return "Person: " + p.Name + ", Email: " + p.Email + ", Phone: " + p.Phone
}

func (p *Person) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Email == "" {
		return errors.New("email is required")
	}
	return nil
}
