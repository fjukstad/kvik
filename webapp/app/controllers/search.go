package controllers

import "github.com/revel/revel"

type Search struct {
	*revel.Controller
}

type Response struct {
	Term string `json:"term"`
}

func (c Search) New(term string) revel.Result {
	res := Response{Term: term}
	return c.RenderJson(res)
}
