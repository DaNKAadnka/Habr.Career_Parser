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

	oldVacancies, err := h.service.Vacancies.GetAllWithFiltration(input)
	if err != nil {
		logrus.Errorf("Error occured while getting old vacancies:\n %s", err.Error())
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Only one service needs to check old vacancies
	if input.IsChosen {
		vacancyForDelete := checkOldVacancies(oldVacancies)
		if len(vacancyForDelete) != 0 {
			err := h.service.Vacancies.DeleteUnactual(vacancyForDelete)
			if err != nil {
				logrus.Errorf("Error occured while deleting old vacancies:\n %s", err.Error())
				newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	vacancies := parse_habr(input.Name, input.Company, input.Salary)

	err = h.service.Vacancies.InsertAll(vacancies)
	if err != nil {
		logrus.Errorf("Error occured while inserting new vacancies:\n %s\n", err.Error())
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}

// check if vacancies in db is actual
// also removing unactual vacancies from oldVacancies
func checkOldVacancies(oldVacancies []parser.Vacancy) []int {

	j := 0
	// id of vacancies that needed to be removed
	var vacancyForDelete []int
	for _, vacancy := range oldVacancies {

		// If resume is not actual
		if resp, err := http.Get(vacancy.Url); err != nil || resp.StatusCode == http.StatusNotFound {
			vacancyForDelete = append(vacancyForDelete, vacancy.Id)
			oldVacancies[j] = vacancy
			j++
		}
	}
	oldVacancies = oldVacancies[:j]
	fmt.Println(oldVacancies)
	return vacancyForDelete
}

func parse_habr(name string, company string, salary int) []parser.Vacancy {

	c := colly.NewCollector()
	var vacancies []parser.Vacancy

	mainLink := "https://career.habr.com/"

	filled_links := make(map[string]bool)
	habrVacanciesLink := fmt.Sprintf("https://career.habr.com/vacancies?type=all&q=%s&salary=%d",
		name+"+"+company, salary) + "&page=%d"

	c.OnHTML(".vacancy-card", func(h *colly.HTMLElement) {
		cardCompany := h.ChildText(".vacancy-card__company-title")
		companyName := h.ChildText(".vacancy-card__skills")
		filled_links[h.Request.URL.String()] = true
		if company == "" || strings.EqualFold(company, cardCompany) {
			salaryString := h.ChildText(".basic-salary")
			minPayment, maxPayment := parseSalary(salaryString)

			vacancies = append(vacancies, parser.Vacancy{
				Url:         mainLink + h.ChildAttr("a.vacancy-card__title-link", "href"),
				Name:        h.ChildText("a.vacancy-card__title-link"),
				MinPayment:  &minPayment,
				MaxPayment:  &maxPayment,
				Description: &companyName,
				Company:     cardCompany,
			})
		}
	})

	c.OnScraped(func(r *colly.Response) {
		page, err := strconv.Atoi(r.Request.URL.Query().Get("page"))
		fmt.Println(r.Request.URL.String())
		if err != nil {
			fmt.Println("Could not convert page param: ", r.Request.URL)
		}
		newLink := fmt.Sprintf(habrVacanciesLink, page+1)
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
