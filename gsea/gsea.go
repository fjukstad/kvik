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

func SystematicName(geneset string) (string, error) {
	return RetrieveFromGSEA(geneset, "Systematic name")
}
func Collection(geneset string) ([]string, error) {
	collection, err := RetrieveFromGSEA(geneset, "Collection")
	if err != nil {
		return []string{}, err
	}
	collections := strings.Split(collection, "      ")
	return collections, nil
}

// Get pubmed id of source publication
func SourcePublication(geneset string) (string, error) {
	pub, err := RetrieveFromGSEA(geneset, "Source publication")
	if err != nil {
		return "", err
	}
	if pub == " " {
		return "", err
	}
	// get pub id
	id := strings.Split(pub, " ")
	if len(pub) > 0 {
		pub = id[1]
	}
	pub = strings.Split(pub, "   ")[0]

	return pub, err
}

func PublicationURL(pubmedID string) string {
	return "https://www.ncbi.nlm.nih.gov/pubmed?CrntRpt=DocSum&cmd=search&term=" + pubmedID
}

func RelatedGeneSets(geneset string) ([]string, error) {
	gs, err := RetrieveFromGSEA(geneset, "Related gene sets")
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(gs, "\n")
	geneSets := []string{}
	for _, line := range lines {
		if strings.Contains(line, "_") {
			line = strings.Trim(line, " ")
			geneSets = append(geneSets, line)
		}
	}
	return geneSets, err
}

func Organism(geneset string) (string, error) {
	return RetrieveFromGSEA(geneset, "Organism")
}

func ContributedBy(geneset string) (string, error) {
	return RetrieveFromGSEA(geneset, "Contributed by")
}

func SourcePlatform(geneset string) (string, error) {
	return RetrieveFromGSEA(geneset, "Source platform")
}

func CompendiumURL(geneset, compendium string) string {
	return "http://software.broadinstitute.org/gsea/msigdb/compendium.jsp?geneSetName=" + geneset + "&compendiumId=" + compendium
}

func DatasetReference(geneset string) ([]string, error) {
	res, err := RetrieveFromGSEA(geneset, "Dataset references")
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(res, "\n")
	geneSets := []string{}
	for _, line := range lines {
		if !strings.Contains(line, "(") {
			line = strings.Trim(line, " ")
			if line != "" {
				geneSets = append(geneSets, line)
			}
		}
	}
	return geneSets, nil
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
