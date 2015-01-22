package models

import (
	"log"
	"strings"

	"github.com/fjukstad/kvik/kegg"
)

type KEGG struct {
	Pathways []kegg.Pathway
	Genes    []kegg.Gene
}

func Init() *KEGG {

	log.Println("Init of kegg model")
	k := new(KEGG)

	k.Pathways = kegg.GetAllHumanPathways()

	return k
}

type Result struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}

func (k KEGG) Search(x string) []Result {

	log.Println(k.Pathways)

	result := SearchPathways(k.Pathways, x)

	return result
}

func SearchPathways(pathways []kegg.Pathway, x string) []Result {
	var results []Result
	for _, pathway := range pathways {
		a := strings.ToLower(pathway.Name)
		b := strings.ToLower(x)
		if strings.Contains(a, b) {
			result := Result{pathway.Id, pathway.Name, pathway.Name}
			results = append(results, result)
		}
	}
	return results
}
