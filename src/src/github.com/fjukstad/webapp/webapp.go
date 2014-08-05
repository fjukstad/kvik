package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/fjukstad/kegg"
)

type InputList struct {
	Input []string
}

type Selection struct {
	Selection []string
}

type Page struct {
	Title string
	Body  []byte
}

var defaultTemplatePaths = []string{
	"templates/_base.html",
	"templates/header.html",
	"templates/navbar.html",
	"templates/footer.html",
}

var emptyTemplatePath = []string{"templates/_empty_page.html",
	"templates/header.html",
}

var indexTemplatePath = "templates/index.html"
var aboutTemplatePath = "templates/about.html"
var browserTemplatePath = "templates/browser.html"
var browserVisualizationTemplatePath = "templates/visualization.html"

var indexTemplate = template.Must(template.ParseFiles(
	append(defaultTemplatePaths, indexTemplatePath)...,
))
var aboutTemplate = template.Must(template.ParseFiles(
	append(defaultTemplatePaths, aboutTemplatePath)...,
))
var browserTemplate = template.Must(template.ParseFiles(
	append(defaultTemplatePaths, browserTemplatePath)...,
))

var browserVisualizationTemplate = template.Must(template.ParseFiles(
	append(defaultTemplatePaths, browserVisualizationTemplatePath)...,
))

func renderTemplate(t *template.Template, w http.ResponseWriter,
	d interface{}) {

	// Cross domain requests in browser are a go go
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Applies template to dataobject d
	err := t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(indexTemplate, w, nil)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(aboutTemplate, w, nil)
}

const lenPath = len("/browser/")

func browserHandler(w http.ResponseWriter, r *http.Request) {
	// Get selected pathways (if any)
	selectedPathways := r.URL.Path[lenPath:]
	log.Print("title of page:", selectedPathways)

	// if user has selected pathways render a visualization
	if len(selectedPathways) > 1 {
		selection := Selection{parsePathwayInput(selectedPathways)}
		renderTemplate(browserVisualizationTemplate, w, selection)
		return
	}

	// if user has not selected any pathways, fetch
	// availible pathways from db and display them to the user
	pathways := kegg.GiveMeSomePathways()
	input := InputList{pathways}

	renderTemplate(browserTemplate, w, input)
}

func parsePathwayInput(input string) []string {
	// Remove any unwanted characters
	a := strings.Replace(input, "%3A", ":", -1)
	a = strings.Replace(a, "&", "", -1)
	a = strings.Replace(a, "=", "", -1)

	// Split into separate hsa:... strings
	b := strings.Split(a, "pathwaySelect")

	// Clear out first empty item
	b = b[1:len(b)]

	return b

}

func parseGeneInput(input string) []string {

	// Remove any unwanted characters
	a := strings.Replace(input, "%3A", ":", -1)
	a = strings.Replace(a, "&", "", -1)
	a = strings.Replace(a, "=", "", -1)

	// Split into separate hsa:... strings
	b := strings.Split(a, "geneSelect")

	// Clear out first empty item
	b = b[1:len(b)]

	return b

}

func main() {

	// cmd line flags
	var ip = flag.String("ip", "localhost", "ip to run on")
	var port = flag.String("port", ":8000", "port to run on")
	flag.Parse()

	address := *ip + *port

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/browser/", browserHandler)

	// Handling requests for css files
	http.Handle("/css/", http.StripPrefix("/css/",
		http.FileServer(http.Dir("css"))))

	// Handling requests for js files
	http.Handle("/js/", http.StripPrefix("/js/",
		http.FileServer(http.Dir("js"))))

	// Handling requests for any public file
	http.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir("public"))))

	// lib directory, d3 jquery and so on.
	http.Handle("/lib/", http.StripPrefix("/lib/",
		http.FileServer(http.Dir("lib"))))

	fmt.Println("Webserver started on", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Panic("Could not start webapp! ", err.Error())
	}

}
