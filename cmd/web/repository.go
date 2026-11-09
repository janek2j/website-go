package main

type Repository interface {
	GetById(id int) any
	GetAll() any
	Create()
	Update()
	Delete()
}

type SQLiteRepository struct {
}

func (r *SQLiteRepository) GetById(id int) any {
	return "not implemented"
}

func (r *SQLiteRepository) GetAll() any {
	return "not implemented"
}

func (r *SQLiteRepository) Create() {
}

func (r *SQLiteRepository) Update() {
}

func (r *SQLiteRepository) Delete() {
}
