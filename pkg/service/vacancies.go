package service

import (
	parser "habr-career"
	"habr-career/pkg/repository"
)

type VacService struct {
	repo repository.Vacancies
}

func NewVacService(repo repository.Vacancies) *VacService {
	return &VacService{repo}
}

func (s *VacService) InsertAll(vacancies []parser.Vacancy) error {
	return s.repo.InsertAll(vacancies)
}
