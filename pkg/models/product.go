package models

import (
	"errors"
	"fmt"
	"time"
)

type ProductCategory struct {
	Name string
}

type Product struct {
	Id                 string
	Name               string
	Decscription       string
	Price              float64
	Category           ProductCategory
	CreatedAt          time.Time
	UpdatedAt          time.Time
	IsAvailableForRent bool
	IsForSale          bool
	IsForExchange      bool
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("name required")
	}
	if p.Price < 0 {
		return errors.New("invalid price")
	}
	return nil
}

func (p *Product) String() string {
	return "Product: " + p.Name + ", Price: " + fmt.Sprintf("%.2f", p.Price) + ", Category: " + p.Category.Name
}
