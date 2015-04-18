package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fjukstad/kvik/kegg"
)

func retrievePathways(filename string) (pathways []string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return pathways, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pathway := scanner.Text()
		pathways = append(pathways, pathway)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return pathways, err
}

func main() {

	var pathwayFilename = flag.String("pathway-file", "pathways",
		"file containing the pathways of interest")

	pathways, err := retrievePathways(*pathwayFilename)
	if err != nil {
		fmt.Println("Error reading pathways from", pathwayFilename, ":", err)
	}

	for _, p := range pathways {

		pathway := kegg.GetPathway(p)
		fmt.Print(pathway.Id + ",")

		genes := pathway.Genes

		for _, g := range genes {
			gene := kegg.GetGene(g)
			fmt.Print(gene.Name + ", ")
		}

		fmt.Print("\n")
	}

	return
}
