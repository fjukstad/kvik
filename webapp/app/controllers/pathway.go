package controllers

import (
	"strings"

	"github.com/fjukstad/kvik/kegg"
	"github.com/revel/revel"
)

type Pathway struct {
	*revel.Controller
}

type Pathways struct {
	Pathways []string `json:"pathways"`
}

func (c Pathway) Index() revel.Result {
	return c.Render()
}

func (c Pathway) Vis(id string) revel.Result {
	ids := strings.Split(id, "+")
	return c.Render(ids)
}

func (c Pathway) JSON(id string) revel.Result {
	graph := kegg.PathwayGraph(id)
	return c.RenderJson(graph)
}

func (c Pathway) Info(id string) revel.Result {
	pw := kegg.GetPathway(id)
	return c.RenderJson(pw.JSON())
}
