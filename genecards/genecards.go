package genecards

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fjukstad/gocache"
)

func Summary(name string) (summary string) {

	url := "http://v4.genecards.org/cgi-bin/carddisp.pl?gene=" + name
	resp, err := gocache.Get(url)
	if err != nil {
		log.Panic(err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Panic(err)
	}

	doc.Find(".gc-subsection").Each(func(i int, s *goquery.Selection) {
		head := s.Find("h3").Text()
		if strings.Contains(head, "Entrez Gene Summary for") {
			summary = s.Find("p").Text()
			return
		}
	})

	return summary

}
