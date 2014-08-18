package main

import (
	"github.com/fjukstad/kegg"
)

func main() {

	geneId := "10458"
	gene := kegg.GetGene(geneId)
	gene.Print()

	pathwayId := "hsa05200"
	pathway := kegg.GetPathway(pathwayId)
	pathway.Print()

	compoundId := "C00575"
	kegg.GetCompound(compoundId)

	return
}
