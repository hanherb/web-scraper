package controllers

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

type fromWeb struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Merchant string `json:"merchant"`
	Image    string `json:"image"`
}

func Scrape() {
	url := "https://www.tokopedia.com/p/handphone-tablet/handphone"

	webData := []fromWeb{}

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("div[class='css-bk6tzz e1nlzfl3']").Each(func(i int, s *goquery.Selection) {
				webData = append(webData, fromWeb{})
				webData[i].URL, _ = s.Find("a").Attr("href")
				webData[i].Name = s.Find("a > div[class='css-16vw0vn'] > div[class='css-11s9vse'] > span").Text()
				webData[i].Price = s.Find("a > div[class='css-16vw0vn'] > div[class='css-11s9vse'] > div > div > span[class='css-o5uqvq']").Text()
				webData[i].Merchant = s.Find("a > div[class='css-16vw0vn'] > div[class='css-11s9vse'] > div[class='css-tpww51'] > div > span:nth-child(2)").Text()
				webData[i].Image, _ = s.Find("a > div[class='css-16vw0vn'] > div > div > div > img").Attr("src")
			})
		},
	}).Start()

	toCSV(webData)
}

func toCSV(m []fromWeb) {
	file, err := os.Create("export.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"url",
		"name",
		"price",
		"merchant",
		"image",
	}

	writer.Write(headers)

	for key := range m {

		r := make([]string, 0, 1+len(headers))

		r = append(
			r,
			m[key].URL,
			m[key].Name,
			m[key].Price,
			m[key].Merchant,
			m[key].Image,
		)

		writer.Write(r)
	}
}
