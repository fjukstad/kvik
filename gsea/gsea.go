package gsea

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fjukstad/gocache"
)

func Abstract(geneset string) (abstract string, err error) {
	url := "http://www.broadinstitute.org/gsea/msigdb/geneset_page.jsp?geneSetName=" + geneset
	resp, err := gocache.Get(url)
	if err != nil {
		fmt.Println("Could not get abstract")
		return "", err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println("could not create goquery doc")
		return "", err
	}

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		rows := s.Find("tr")
		rows.Each(func(j int, se *goquery.Selection) {
			header := se.Find("th").Text()
			if strings.Contains(header, "Full description") {
				abstract = se.Find("td").Text()
				return
			}
		})
		return
	})
	return abstract, nil
}
