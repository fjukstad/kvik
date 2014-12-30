package main

import (
	"fmt"

	"github.com/fjukstad/kegg"
)

func main() {

	geneId := "10458"
	gene := kegg.GetGene(geneId)

	pathwayId := "hsa05200"
	pathway := kegg.GetPathway(pathwayId)

	compoundId := "C00575"
	kegg.GetCompound(compoundId)

	pws := kegg.GetAllHumanPathways()
	fmt.Println(pws)

	/*
		for _, id := range pws {

			pw := kegg.NewKeggPathway(id)
			log.Println(pw.Name, len(pw.Entries))
		}
	*/

	return
}
