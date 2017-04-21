package db

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fjukstad/kvik/eutils"
)

type Service struct {
	Address string
	Cache   bool
}

func Init(address string, cache bool) Service {
	return Service{address, cache}
}

func (s Service) Summary(geneSymbol string) (string, error) {
	summary, err := s.Doc(geneSymbol)
	if err != nil {
		return "", err
	}
	return summary.Summary, nil
}

func (s Service) Doc(geneSymbol string) (eutils.DocumentSummary, error) {
	u := s.Address + "/doc?geneSymbol=" + geneSymbol
	resp, err := http.Get(u)
	if err != nil {
		return eutils.DocumentSummary{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return eutils.DocumentSummary{}, err
	}

	var summary eutils.DocumentSummary
	err = json.Unmarshal(body, &summary)
	if err != nil {
		return eutils.DocumentSummary{}, err
	}

	return summary, nil
}
