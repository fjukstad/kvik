package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/fjukstad/kvik/genecards"
	"github.com/fjukstad/kvik/kegg"
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

func (c Gene) Fc(genes []string) revel.Result {
	return c.Render()
}

func (c Gene) Exprs(id []string) revel.Result {
	things := dataset.Exprs([]string{"BRCA1", "BRCA2"})
	fmt.Println(things)
	return c.Render()
}

var dataset *models.Dataset

func InitDataset() {
	var err error
	dataset, err = models.InitDataset()
	if err != nil {
		log.Panic(err)
	}
}
