package controllers

import (
	"github.com/fjukstad/kvik/webapp/app/models"
	"github.com/revel/revel"
)

var kegg *models.KEGG

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
	result := kegg.Search(term)
	response := Response{Terms: result}
	return c.RenderJson(response)
}

func InitKEGG() {
	kegg = models.Init()
}
