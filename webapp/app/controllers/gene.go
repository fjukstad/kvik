package controllers

import (
	"strings"

	"github.com/fjukstad/kvik/genecards"
	"github.com/fjukstad/kvik/kegg"
	"github.com/revel/revel"
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
