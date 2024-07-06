package parser

type Vacancy struct {
	Id          int     `json:"id" db:"id"`
	Url         string  `json:"url" db:"url"`
	Name        string  `json:"name" db:"name"`
	MinPayment  *int    `json:"min_payment" db:"min_payment"`
	MaxPayment  *int    `json:"max_payment" db:"max_payment"`
	Description *string `json:"description" db:"description"`
	Company     string  `json:"company" db:"employer_name"`
}

type SearchVacancies struct {
	Name     string `json:"name" validate:"optional"`
	Company  string `json:"company" validate:"optional"`
	Salary   int    `json:"salary" validate:"optional"`
	IsChosen bool   `json:"is_chosen" validate:"required"`
}

type VacancyAnalitics struct {
	Count              int    `json:"count"`
	MostCompany        string `json:"most_company"`
	CountOfMostCompany int    `json:"count_of_most_company"`
	AvaragePayment     int    `json:"avarage_payment"`
}
