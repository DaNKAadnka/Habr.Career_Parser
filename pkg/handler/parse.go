package handler

import (
	"fmt"
	parser "habr-career"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// @Summury ParseVacancies
// @Tags Parse
// @ID parse-v
// @Param input body parser.SearchVacancies true "input"
// @Success 200 {integer} 1
// @Router /vacancies [post]
func (h *Handler) parseVacancies(ctx *gin.Context) {

	var input parser.SearchVacancies

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	vacancies := parse_habr(input.Name, input.Company, input.Salary)

	err := h.service.Vacancies.InsertAll(vacancies)

	if err != nil {
		logrus.Errorf("Vacancies insertion error: %s\n", err.Error())
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}

func parse_habr(name string, company string, salary int) []parser.Vacancy {

	c := colly.NewCollector()
	var vacancies []parser.Vacancy

	filled_links := make(map[string]bool)
	habrVacanciesLink := fmt.Sprintf("https://career.habr.com/vacancies?type=all&q=%s&salary=%d",
		name+"+"+company, salary) + "&page=%d"

	c.OnHTML(".vacancy-card", func(h *colly.HTMLElement) {
		card_company := h.ChildText(".vacancy-card__company-title")
		fmt.Println("Company: ", card_company)
		filled_links[h.Request.URL.String()] = true
		if company == "" || strings.EqualFold(company, card_company) {
			salaryString := h.ChildText(".basic-salary")
			minPayment, maxPayment := parseSalary(salaryString)

			vacancies = append(vacancies, parser.Vacancy{
				Url:         h.ChildAttr("a.vacancy-card__title-link", "href"),
				Name:        h.ChildText("a.vacancy-card__title-link"),
				MinPayment:  minPayment,
				MaxPayment:  maxPayment,
				Description: h.ChildText(".vacancy-card__skills"),
				Company:     card_company,
			})
			fmt.Printf("Appended")
		}
	})

	c.OnScraped(func(r *colly.Response) {
		page, err := strconv.Atoi(r.Request.URL.Query().Get("page"))
		fmt.Println(r.Request.URL.String())
		if err != nil {
			fmt.Println("Could not convert page param: ", r.Request.URL)
		}
		newLink := fmt.Sprintf(habrVacanciesLink, page+1)
		fmt.Println("New link- ", newLink)
		if filled_links[r.Request.URL.String()] {
			c.Visit(newLink)
		}
	})

	c.Visit(fmt.Sprintf(habrVacanciesLink, 1))

	return vacancies

}

func parseSalary(s string) (int, int) {
	min, max := 0, 0
	var err error
	strings := strings.Split(s, " ")
	for i, x := range strings {
		if x == "от" {
			min, err = strconv.Atoi(strings[i+1])
			if err != nil {
				min = 0
			}
		}
		if x == "до" {
			max, err = strconv.Atoi(strings[i+1])
			if err != nil {
				max = 0
			}
		}
	}
	return min, max
}
