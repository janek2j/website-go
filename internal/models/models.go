package models

import (
	"errors"
	"fmt"
	"time"
)

// User is the user model
type User struct {
	Id        uint
	Name      string
	Email     string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Book is book model
type Book struct {
	Id             uint
	ISBN           string
	Title          string
	Authors        string
	Publisher      string
	Year           uint16
	Pages          uint16
	Description    string
	OwnerId        int
	Tags           string
	Pictures       string
	EstimatedValue float32
	CategoryId     int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsPublic       bool
}

type GenericItem struct {
}

// Borrowing is the borrowing model
type Borrowing struct {
	ProductId    uint      `json:"product_id"`
	LenderId     uint      `json:"lender_id"`
	BorrowerId   uint      `json:"borrower_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Picture      string    `json:"picture"`
	IsHanded     bool
	IsReturned   bool
	BorrowerNote string
	LenderNote   string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Book Book
}

// ItemCategory is item category model
type ItemCategory struct {
	Id   uint
	Name string
}

type LendingOffer struct {
}

type BorrowingRequest struct {
}

// DescriptionTemplate is a descripton template model
type DescriptionTemplate struct {
	Id       uint
	Template string
}

func (p *User) String() string {
	return "User: " + p.Name + ", Email: " + p.Email
}

func (p *User) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

//-------------------------------------

// Reservation holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	BookId    uint
}

type Product struct {
	Id                 string
	Name               string
	Decscription       string
	Price              float64
	Category           ItemCategory
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
