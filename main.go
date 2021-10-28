package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type section struct {
	Number      int
	Name        string
	SectionType string
	Status      string
}

type course struct {
	Subject  string
	Number   int
	Title    string
	Sections []section
}

func main() {
	subjectCollector := colly.NewCollector()
	subjects := []string{}
	courses := []course{}

	subjectCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	subjectCollector.OnHTML(".browse-by-subjects .browse-subjectcode a", func(e *colly.HTMLElement) {
		fmt.Printf("Found subject:\t%s\n", e.Text)
		subjects = append(subjects, e.Text)
	})

	subjectCollector.Visit("https://classes.cornell.edu")

	c := colly.NewCollector()

	c.OnHTML("[aria-label^=\"Course\"]", func(e *colly.HTMLElement) {
		catalogNum, err := strconv.Atoi(e.Attr("data-catalog-nbr"))
		if err != nil {
			fmt.Printf("Error parsing catalog number: %v\nContinuing...\n", err)
			return
		}

		crs := course{
			Subject: e.Attr("data-subject"),
			Number:  catalogNum,
			Title:   e.DOM.Find(".title-coursedescr").Text(),
		}

		sections := []section{}

		e.DOM.Find("ul").Each(func(i int, e *goquery.Selection) {
			classNumbers := e.Find(".class-numbers > p").Text()
			if len(classNumbers) == 0 {
				return
			}
			numberNameLst := strings.Split(classNumbers, crs.Subject+" "+strconv.Itoa(crs.Number))
			sectionNumber, err := strconv.Atoi(strings.TrimSpace(numberNameLst[0]))
			if err != nil {
				fmt.Printf("Error parsing section number: %v\nContinuing...\n", err)
				return
			}
			sectionName := strings.TrimSpace(numberNameLst[1])
			sectionType := strings.TrimSpace(strings.Split(numberNameLst[1], " ")[0])
			sectionStatus := e.Find(".open-status > span").AttrOr("data-content", "Unknown")

			s := section{
				Number:      sectionNumber,
				Name:        sectionName,
				SectionType: sectionType,
				Status:      sectionStatus,
			}
			sections = append(sections, s)
		})

		crs.Sections = sections
		courses = append(courses, crs)
	})

	subjectCollector.Wait()

	for _, v := range subjects {
		c.Visit("https://classes.cornell.edu/browse/roster/SP21/subject/" + v)
		fmt.Println("Visiting", v)
		time.Sleep(time.Second)
	}

	c.Wait()

	output, err := json.Marshal(courses)
	if err != nil {
		panic(err)
	}

	_ = ioutil.WriteFile("roster-data.json", output, 0644)
}
