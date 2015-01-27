package controllers

import (
	"github.com/fjukstad/kvik/kegg"
	"github.com/revel/revel"
)

type Pathway struct {
	*revel.Controller
}

func (c Pathway) Index() revel.Result {
	return c.Render()
}

func (c Pathway) Vis(id string) revel.Result {
	return c.Render()
}

func (c Pathway) JSON(id string) revel.Result {
	graph := kegg.PathwayGraph(id)
	return c.RenderJson(graph)
}
