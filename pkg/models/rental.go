package models

import "time"

type Rental struct {
	ProductId  string    `json:"product_id"`
	LenderId   string    `json:"lender_id"`
	BorrowerId string    `json:"borrower_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Picture    string    `json:"picture"`
}
