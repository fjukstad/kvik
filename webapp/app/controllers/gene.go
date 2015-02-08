package controllers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/fjukstad/kvik/genecards"
	"github.com/fjukstad/kvik/kegg"
	"github.com/fjukstad/kvik/utils"
	"github.com/revel/revel"

	"github.com/fjukstad/kvik/webapp/app/models"
)

type Gene struct {
	*revel.Controller
}

func (c Gene) Info(id string) revel.Result {
	g := kegg.GetGene(id)

	// exchange definition with gene cards summary
	name := strings.Split(g.Name, " ")[0]
	name = strings.TrimRight(name, ",")
	g.Definition = genecards.Summary(name)

	return c.RenderJson(g.JSON())
}

func (c Gene) Fc() revel.Result {
	genes := getGenes(c)
	geneNames := idToNames(genes)
	output := dataset.Fc(geneNames)
	result, err := prepareResult(genes, output)
	if err != nil {
		c.RenderText(fmt.Sprintf("%s", err))
	}
	response := utils.ClientCompResponse{result}
	return c.RenderJson(response)
}

func (c Gene) Exprs() revel.Result {
	genes := getGenes(c)
	geneNames := idToNames(genes)
	output := dataset.Exprs(geneNames)
	result, err := prepareResult(genes, output)
	if err != nil {
		c.RenderText(fmt.Sprintf("%s", err))
	}
	response := utils.ClientCompResponse{result}
	return c.RenderJson(response)
}

func prepareResult(genes, output []string) (map[string]string, error) {
	if len(genes) != len(output) {
		return nil, errors.New("List of genes and output from data is of different lengths" + strconv.Itoa(len(genes)) + " and " + strconv.Itoa(len(output)))
	}
	result := make(map[string]string)
	for i, g := range genes {
		result[g] = output[i]
	}
	return result, nil
}

func getGenes(c Gene) []string {
	genes := c.Params.Get("genes")
	return strings.Split(genes, "+")
}

func idToNames(ids []string) []string {
	var names []string
	for _, id := range ids {
		gene := kegg.GetGene(id)
		name := strings.Split(gene.Name, " ")[0]
		names = append(names, name)
	}
	return names
}

var dataset *models.Dataset

func InitDataset() {
	var err error
	dataset, err = models.InitDataset()
	if err != nil {
		log.Panic(err)
	}
}
