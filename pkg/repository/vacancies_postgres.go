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
		(:url, :name, :min_payment, :max_payment, :description, :employer_name)
		ON CONFLICT DO NOTHING`

	_, err := r.db.NamedExec(query, vacancies)
	return err
}

func (r *VacPostgres) GetAllWithFiltration(filters parser.SearchVacancies) ([]parser.Vacancy, error) {

	var query strings.Builder
	var oldVacancies []parser.Vacancy

	query.WriteString(`SELECT * FROM vacancies`)

	and_flag := false

	if filters.Name != "" {
		likeName := "%" + filters.Name + "%"
		query.WriteString(fmt.Sprintf(" WHERE (name LIKE '%s') OR (description LIKE '%s')", likeName, likeName))
		and_flag = true
	}
	if filters.Company != "" {
		if and_flag {
			query.WriteString(" AND")
		} else {
			query.WriteString(" WHERE")
		}
		query.WriteString(fmt.Sprintf(" employer_name LIKE '%s%s%s'", "%", filters.Company, "%"))
	}
	if filters.Salary != 0 {
		if and_flag {
			query.WriteString(" AND")
		} else {
			query.WriteString(" WHERE")
		}
		query.WriteString(fmt.Sprintf(` min_payment < %d`, filters.Salary))
	}

	if err := r.db.Select(&oldVacancies, query.String()); err != nil {
		return nil, err
	}

	return oldVacancies, nil
}

func (r *VacPostgres) DeleteUnactual(ids []int) error {

	var query strings.Builder
	query.WriteString(`DELETE FROM vacancies WHERE id IN (`)
	for i, id := range ids {
		query.WriteString(fmt.Sprintf("%d", id))
		if i != len(ids)-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(")")
	fmt.Println(query.String())

	_, err := r.db.Exec(query.String())
	return err
}
