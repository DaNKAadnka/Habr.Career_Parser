package parser

type Vacancy struct {
	Id          int    `json:"id"`
	Url         string `json:"url"`
	Name        string `json:"name"`
	MinPayment  int    `json:"min_payment"`
	MaxPayment  int    `json:"max_payment"`
	Description string `json:"description"`
	Company     string `json:"company"`
}

type SearchVacancies struct {
	Name    string `json:"name" validate:"optional"`
	Company string `json:"company" validate:"optional"`
	Salary  int    `json:"salary" validate:"optional"`
}
