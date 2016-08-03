package gsea

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fjukstad/gocache"
)

func Abstract(geneset string) (abstract string, err error) {
	return RetrieveFromGSEA(geneset, "Full description")
}

func BriefDescription(geneset string) (briefdescription string, err error) {
	return RetrieveFromGSEA(geneset, "Brief description")
}

func RetrieveFromGSEA(geneset, info string) (fromGSEA string, err error) {
	URL := "http://www.broadinstitute.org/gsea/msigdb/geneset_page.jsp?geneSetName=" + geneset
	resp, err := gocache.Get(URL)
	if err != nil {
		fmt.Println("Could not get abstract")
		return "", err
	}

	//resp.Request.URL = url
	resp.Request = &http.Request{}
	u, _ := url.Parse(URL)
	resp.Request.URL = u

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println("could not create goquery doc")
		return "", err
	}

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		rows := s.Find("tr")
		rows.Each(func(j int, se *goquery.Selection) {
			header := se.Find("th").Text()
			if strings.Contains(header, info) {
				fromGSEA = se.Find("td").Text()
				return
			}
		})
		return
	})
	return fromGSEA, nil
}
