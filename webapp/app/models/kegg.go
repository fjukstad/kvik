package models

import (
	"strings"

	"github.com/fjukstad/kvik/kegg"
)

type KEGG struct {
	Pathways []kegg.Pathway
	Genes    []kegg.Gene
}

func Init() *KEGG {
	k := new(KEGG)
	k.Pathways = kegg.GetAllHumanPathways()
	k.Genes = kegg.GetAllHumanGenes()
	return k
}

type Result struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}

func (k KEGG) Search(x string) []Result {
	result := SearchPathways(k.Pathways, x)
	result = append(result, SearchGenes(k.Genes, x)...)
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

func SearchGenes(genes []kegg.Gene, x string) []Result {
	var results []Result
	for _, gene := range genes {
		a := strings.ToLower(gene.Name)
		b := strings.ToLower(x)
		if strings.Contains(a, b) {
			result := Result{gene.Id, gene.Name, gene.Name}
			results = append(results, result)
		}
	}
	return results
}
