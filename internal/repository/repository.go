package repository

import "github.com/janek2j/website-go/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error

	InsertBook(b models.Book) error
}
