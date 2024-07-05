package repository

import (
	parser "habr-career"

	"github.com/jmoiron/sqlx"
)

type Vacancies interface {
	InsertAll(vacancies []parser.Vacancy) error
	GetAllWithFiltration(filters parser.SearchVacancies) ([]parser.Vacancy, error)
}

type Resume interface {
}

type Repository struct {
	Vacancies
	Resume
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Vacancies: NewVacPostgres(db),
	}
}
