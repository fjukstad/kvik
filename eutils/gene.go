package eutils

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GenomicInfo struct {
	ChrLoc    int
	ChrAccVer string
	ChrStart  int
	ChrEnd    int
	ExonCount int
}

type Organism struct {
	ScientificName string
	CommonName     string
	TaxId          int
}

type LocationHistType struct {
	AnnotationRelease int
	AssemblyAccVer    string
	ChrAccVer         string
	ChrStart          int
	ChrEnd            int
}

var geneEndpoint = "http://eutils.ncbi.nlm.nih.gov/entrez/eutils/esummary.fcgi?db=gene&id="

func GetDocumentSummary(id string) (*DocumentSummary, error) {
	u := geneEndpoint + id
	resp, err := http.Get(u)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SummaryResult
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.DocumentSummarySet.Status != "OK" {
		return nil, errors.New("Could not get information about gene with entrez id " + id)
	}

	return &result.DocumentSummarySet.DocumentSummary[0], nil

}
