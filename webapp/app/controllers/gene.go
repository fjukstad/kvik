package controllers

import (
	"github.com/fjukstad/kvik/kegg"
	"github.com/revel/revel"
)

type Gene struct {
	*revel.Controller
}

func (c Gene) Info(id string) revel.Result {
	g := kegg.GetGene(id)
	return c.RenderJson(g.JSON())
}
