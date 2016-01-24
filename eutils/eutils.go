package eutils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type SearchResult struct {
	Count    int
	RetMax   int
	RetStart int
	IdList   []int `xml:"IdList>Id"`
	/*
		TranslationSet   []Translation
		TranslationStack []TermSet
		QueryTranslation string
	*/
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

	resp, err := http.Get(url)

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

type SummaryResult struct {
	Status             string `xml:"status,attr"`
	DocumentSummarySet DocumentSummarySet
}

type DocumentSummarySet struct {
	DbBuild         string
	DocumentSummary []DocumentSummary
}

type DocumentSummary struct {
	Uid             string `xml:"uid,attr"`
	PubDate         string
	Source          string
	Authors         []Authors
	Title           string
	SortTitle       string
	Volume          string
	Issue           string
	Pages           string
	Lang            string `xml:"Lang>string"`
	NlmUniqueID     int
	ISSN            string
	ESSN            string
	PubType         string `xml:"PubType>flag"`
	RecordStatus    string
	PubStatus       int
	ArticleIds      []ArticleId
	History         []PubMedPubDate
	Attributes      string `xml:"Attributes>flag"`
	PmcRefCount     int
	FullJournalName string
	ViewCount       int
	DocType         string
	SortPubDate     string
	SortFirstAuthor string
}

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

	resp, err := http.Get(url)

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
