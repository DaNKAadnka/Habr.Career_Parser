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
// @Success 200 {object} parser.VacancyAnalitics
// @Router /vacancies [post]
func (h *Handler) parseVacancies(ctx *gin.Context) {

	var input parser.SearchVacancies

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	newVacancies := parse_habr(input.Name, input.Company, input.Salary)

	if len(newVacancies) != 0 {

		err := h.service.Vacancies.InsertAll(newVacancies)
		if err != nil {
			logrus.Errorf("Error occured while inserting new vacancies:\n %s\n", err.Error())
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	vacancies, err := h.service.Vacancies.GetAllWithFiltration(input)
	if err != nil {
		logrus.Errorf("Error occured while getting old vacancies:\n %s", err.Error())
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Only one service needs to check old vacancies
	if input.IsChosen {
		var vacancyForDelete []int
		vacancyForDelete, vacancies = checkOldVacancies(vacancies)
		if len(vacancyForDelete) != 0 {
			err := h.service.Vacancies.DeleteUnactual(vacancyForDelete)
			if err != nil {
				logrus.Errorf("Error occured while deleting old vacancies:\n %s", err.Error())
				newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}
		fmt.Println("Delete Done")
	}

	for i, s := range vacancies {
		fmt.Println(i, s)
	}
	response := calculateAnalysis(vacancies)

	fmt.Println(response)
	// Calculate statistics
	ctx.JSON(http.StatusOK, response)
}

// check if vacancies in db is actual
// also removing unactual vacancies from oldVacancies
func checkOldVacancies(vacancies []parser.Vacancy) ([]int, []parser.Vacancy) {

	j := 0
	// id of vacancies that needed to be removed
	var vacancyForDelete []int
	for _, vacancy := range vacancies {

		// If resume is not actual
		if resp, err := http.Get(vacancy.Url); err != nil || resp.StatusCode == http.StatusNotFound {
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}
			vacancyForDelete = append(vacancyForDelete, vacancy.Id)
			vacancies[j] = vacancy
			j++
		}
	}
	vacancies = vacancies[j:]
	// for _, s := range vacancyForDelete {
	// 	fmt.Println("Old: ", s)
	// }
	return vacancyForDelete, vacancies
}

func parse_habr(name string, company string, salary int) []parser.Vacancy {

	c := colly.NewCollector()
	var vacancies []parser.Vacancy

	mainLink := "https://career.habr.com"

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
		// fmt.Println(r.Request.URL.String())
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

func calculateAnalysis(vacancies []parser.Vacancy) parser.VacancyAnalitics {
	var analitics parser.VacancyAnalitics
	analitics.Count = len(vacancies)

	paymentSum := 0
	paymentCnt := 0
	companyCount := make(map[string]int)
	for _, vacancy := range vacancies {
		companyCount[vacancy.Company] += 1
		currentAvarage := 0
		if vacancy.MinPayment != nil {
			currentAvarage += *vacancy.MinPayment
		}
		if vacancy.MaxPayment != nil {
			currentAvarage += *vacancy.MaxPayment
		}
		if vacancy.MinPayment != nil && vacancy.MaxPayment != nil {
			currentAvarage /= 2
		}

		if currentAvarage != 0 {
			paymentCnt += 1
		}
		paymentSum += currentAvarage
	}
	// Avarage Calculated
	if paymentCnt == 0 {
		analitics.AvaragePayment = 0
	} else {
		analitics.AvaragePayment = paymentSum / paymentCnt
	}

	maxCompany := ""
	maxCompanyCnt := 0
	for name, cnt := range companyCount {
		if cnt > maxCompanyCnt {
			maxCompanyCnt = cnt
			maxCompany = name
		}
	}
	analitics.MostCompany = maxCompany
	analitics.CountOfMostCompany = maxCompanyCnt

	return analitics
}
