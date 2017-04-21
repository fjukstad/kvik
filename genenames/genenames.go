package genenames

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fjukstad/gocache"
	"github.com/pkg/errors"
)

type Result struct {
	ResponseHeader `json:"responseHeader"`
	Response       `json:"response"`
}

type ResponseHeader struct {
	Status int `json:"status"`
	QTime  int
}

type Response struct {
	NumFound int   `json:"numFound"`
	Start    int   `json:"stat"`
	Docs     []Doc `json:"docs"`
}

type Doc struct {
	HGNCId   string `json:"hgnc_id"`
	Symbol   string `json:"symbol"`
	EntrezId string `json:"entrez_id"`
}

var baseUrl = "http://rest.genenames.org/"

// Fetches information about a gene with the given gene symbol.
// Warning: fetch might return more than one result, but we're returning
// only the first.
func GetDoc(symbol string) (Doc, error) {
	u := baseUrl + "fetch/symbol/" + symbol

	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("Accept", "application/json")

	resp, err := gocache.Do(req)
	if err != nil {
		return Doc{}, errors.Wrap(err, "Could not download information about gene "+symbol)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Doc{}, errors.Wrap(err, "Could not read response body for gene "+symbol)
	}

	var result Result
	err = json.Unmarshal(body, &result)

	if err != nil {
		return Doc{}, errors.Wrap(err, "Could not unmarshal json from genenames: "+string(body))
	}

	if len(result.Response.Docs) == 0 {
		return Doc{}, errors.New("Could not get doc for gene with symbol " + symbol)
	}

	return result.Response.Docs[0], nil

}
