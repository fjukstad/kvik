package controllers

import (
	"github.com/fjukstad/kvik/webapp/app/models"
	"github.com/revel/revel"
)

var Kegg *models.KEGG

type Search struct {
	*revel.Controller
}

type Response struct {
	Terms []models.Result
}

type Content struct {
	Pathways map[string]string
	Genes    map[string]string
}

func (c Search) New(term string) revel.Result {
	result := Kegg.Search(term)
	response := Response{Terms: result}
	return c.RenderJson(response)
}

func InitKEGG() {
	Kegg = models.Init()
}

func (c Search) Pathways() revel.Result {
	result := Kegg.Search("")
	response := Response{Terms: result}
	return c.RenderJson(response)
}
