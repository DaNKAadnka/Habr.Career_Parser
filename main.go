package main

import (
	"fmt"

	"strconv"

	"github.com/gocolly/colly"
)

func main() {

	c := colly.NewCollector()
	var vacanciesLinks []string

	i := 1
	cnt := 0
	habrVacanciesLink := "https://career.habr.com/vacancies?type=all&q=Backend&page="

	c.OnHTML(".vacancy-card__title", func(h *colly.HTMLElement) {
		newLink := h.ChildAttr("a", "href")
		// replace with set structure
		cnt += 1
		vacanciesLinks = append(vacanciesLinks, newLink)
	})

	c.OnScraped(func(r *colly.Response) {
		i += 1
		fmt.Println(r.Request.URL, cnt)
		if cnt > 0 {
			cnt = 0
			c.Visit(habrVacanciesLink + strconv.Itoa(i))
		}
	})

	c.Visit(habrVacanciesLink + "1")

	for i, s := range vacanciesLinks {
		fmt.Println(i, s)
	}

}
