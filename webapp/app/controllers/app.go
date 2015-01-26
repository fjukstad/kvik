package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	pathways := Kegg.Pathways
	genes := Kegg.Genes
	return c.Render(pathways, genes)
}

func (c App) Pathways() revel.Result {

	return c.Render()
}
