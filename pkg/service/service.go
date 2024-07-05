package service

import (
	parser "habr-career"
	"habr-career/pkg/repository"
)

type Vacancies interface {
	InsertAll(vacancies []parser.Vacancy) error
	GetAllWithFiltration(filters parser.SearchVacancies) ([]parser.Vacancy, error)
	DeleteUnactual(ids []int) error
}

type Resume interface {
}

type Service struct {
	Vacancies
	Resume
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Vacancies: NewVacService(repo.Vacancies),
	}
}
