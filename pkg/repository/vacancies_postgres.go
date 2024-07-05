package repository

import (
	"fmt"
	parser "habr-career"
	"strings"

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

func (r *VacPostgres) GetAllWithFiltration(filters parser.SearchVacancies) ([]parser.Vacancy, error) {

	var query strings.Builder
	var oldVacancies []parser.Vacancy

	query.WriteString(`SELECT (url, name, min_payment, max_payment, employer_name)
		FROM vacancies`)

	and_flag := false

	if filters.Name != "" {
		query.WriteString(fmt.Sprintf(` WHERE name LIKE \%%s\%`, filters.Name))
		and_flag = true
	}
	if filters.Company != "" {
		if and_flag {
			query.WriteString(" AND")
		}
		query.WriteString(fmt.Sprintf(` WHERE employer_name LIKE \%%s\%`, filters.Company))
	}
	if filters.Salary != 0 {
		if and_flag {
			query.WriteString(" AND")
		}
		query.WriteString(fmt.Sprintf(` WHERE min_payment < %d`, filters.Salary))
	}

	if err := r.db.Select(&oldVacancies, query.String()); err != nil {
		return nil, err
	}

	return oldVacancies, nil
}
