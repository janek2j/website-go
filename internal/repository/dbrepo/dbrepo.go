package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/janek2j/website-go/internal/config"
	"github.com/janek2j/website-go/internal/models"
	"github.com/janek2j/website-go/internal/repository"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {

	return &postgresDbRepo{
		App: a,
		DB:  conn,
	}
}

func (m *postgresDbRepo) InsertReservation(res models.Reservation) error {
	return nil
}

// InsertBook inserts a book into the database
func (m *postgresDbRepo) InsertBook(bk models.Book) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into books (title, authors, year, pages, created_at, updated_at) 
			values($1, $2, $3, $4, $5, $6)`
	_, err := m.DB.ExecContext(ctx, stmt,
		bk.Title,
		bk.Authors,
		bk.Year,
		bk.Pages,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
