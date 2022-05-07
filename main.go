package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Fish struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Year   string `json:"year"`
	Region string `json:"region"`
}

func Run() {
	regex := regexp.MustCompile(`\s+`)
	var fishes []Fish
	c := colly.NewCollector()

	c.OnHTML("div.species-directory__species--8col", func(element *colly.HTMLElement) {
		dom := element.DOM
		fish := Fish{
			Name:   dom.Find("div.species-directory__species-title--name").Text(),
			Status: regex.ReplaceAllString(strings.TrimSpace(dom.Find("div.species-directory__species-status-row").Find("div.species-directory__species-status").Text()), " "),
			Year:   regex.ReplaceAllString(strings.TrimSpace(dom.Find("div.species-directory__species-status-row").Find("div.species-directory__species-year").Text()), " "),
			Region: regex.ReplaceAllString(strings.TrimSpace(dom.Find("div.species-directory__species-status-row").Find("div.species-directory__species-region").Text()), " "),
		}
		fishes = append(fishes, fish)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	//
	for i := 1; i <= 5; i++ {
		url := fmt.Sprintf("https://www.fisheries.noaa.gov/species-directory/threatened-endangered?title=&species_category=any&species_status=any&regions=all&items_per_page=25&page=%d&sort=", i)
		err := c.Visit(url)
		if err != nil {
			log.Fatal("Scraper Error!")
			return
		}
	}
	WriteToJsonFile(fishes)
}
func WriteToJsonFile(data []Fish) {
	dataBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal("MarshalIndent Error!")
		return
	}
	_ = ioutil.WriteFile("fishes.json", dataBytes, 0777)

}
func main() {
	Run()
}
