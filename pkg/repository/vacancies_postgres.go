package repository

import (
	parser "habr-career"

	"github.com/jmoiron/sqlx"
)

type VacPostgres struct {
	db *sqlx.DB
}

func NewVacPostgres(db *sqlx.DB) *VacPostgres {
	return &VacPostgres{db}
}

func (r *VacPostgres) InsertAll(vacancies []parser.Vacancy) error {

	query := `INSERT INTO vacancies (url, name, min_payment,
		max_payment, description, employer_name) VALUES
		(:url, :name, :minpayment, :maxpayment, :description, :company)
		ON CONFLICT DO NOTHING`

	_, err := r.db.NamedExec(query, vacancies)
	return err
}
