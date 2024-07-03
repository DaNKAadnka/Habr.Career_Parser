package parser

type Vacancy struct {
	Id            int    `json:"id"`
	Url           string `json:"url"`
	Name          string `json:"name"`
	Min_payment   int    `json:"min_payment"`
	Max_payment   int    `json:"max_payment"`
	Description   string `json:"description"`
	Is_active     bool   `json:"is_active"`
	Employer_name string `json:"employer_name"`
}
