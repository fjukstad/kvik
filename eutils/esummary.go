package eutils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/fjukstad/gocache"
)

type PubMedPubDate struct {
	PubStatus string
	Date      string
}

type ArticleId struct {
	IdType  string
	IdTypeN int
	Value   int
}

type Authors struct {
	Authors    []Author
	LastAuthor string
}

type Author struct {
	Name     string
	AuthType string
}

func Summary(r *SearchResult) (*SummaryResult, error) {

	baseUrl := "http://eutils.ncbi.nlm.nih.gov/entrez/eutils/esummary.fcgi?db=PubMed&id="

	ids := []string{}
	for _, id := range r.IdList {
		ids = append(ids, strconv.Itoa(id))
	}

	url := baseUrl + strings.Join(ids, ",") + "&version=2.0"

	fmt.Println(baseUrl, ids, url)

	resp, err := gocache.Get(url)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := new(SummaryResult)

	err = xml.Unmarshal(body, s)

	return s, nil

}
