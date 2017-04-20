package eutils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fjukstad/gocache"
)

type SearchResult struct {
	Count    int
	RetMax   int
	RetStart int
	IdList   []int `xml:"IdList>Id"`
}

type Id struct {
	Id string `xml:",int"`
}

type Translation struct {
	From string
	To   string
}

type TermSet struct {
	Term    string
	Field   string
	Count   int
	Explode string
	OP      string
}

func Search(terms []string) (*SearchResult, error) {

	baseUrl := "http://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term=science[journal]+AND+"

	searchTerms := strings.Join(terms, "+OR+")

	url := baseUrl + searchTerms + "+AND+breast+cancer"

	resp, err := gocache.Get(url)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := new(SearchResult)

	err = xml.Unmarshal(body, s)

	return s, err
}
