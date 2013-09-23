package main

import (
    "fmt"
    "net/http" 
    "html/template"
    "log"
    "nowac/kegg"
    "strings"
    "flag"
)

type Dataset struct {
    Genes []string
    Exprs map[string] []int 
}

type Page struct {
    Title string
    Body []byte
}

type Selection struct {
    Genes [] string
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
var demoTemplatePath = "templates/demo.html"
var demoVisualizationTemplatePath = "templates/visualization.html"

var indexTemplate = template.Must(template.ParseFiles(
    append(defaultTemplatePaths, indexTemplatePath)...
))
var aboutTemplate = template.Must(template.ParseFiles(
    append(defaultTemplatePaths, aboutTemplatePath)...
))
var demoTemplate = template.Must(template.ParseFiles(
    append(defaultTemplatePaths, demoTemplatePath)...
))

var demoVisualizationTemplate = template.Must(template.ParseFiles(
    append(emptyTemplatePath, demoVisualizationTemplatePath)...
))


var keggInterface kegg.Kegg

func renderTemplate (t *template.Template, w http.ResponseWriter, 
                                                d interface{}){
    
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

const lenPath = len("/demo/")

func demoHandler(w http.ResponseWriter, r *http.Request) {

    selectedGenes := r.URL.Path[lenPath:]
    fmt.Println("title of page: ", selectedGenes)  

    // If user has selected genes, display them
    if len(selectedGenes) > 1 {
        fmt.Println(parseGeneInput(selectedGenes))
        formattedGenes := Selection{parseGeneInput(selectedGenes)}

        renderTemplate(demoVisualizationTemplate, w, formattedGenes)
        return
    }
    fmt.Println(w)
    fmt.Println(r) 
    genes := keggInterface.GetNFirstGeneIDs(100)
    ds := Dataset{genes, nil}
    renderTemplate(demoTemplate, w, ds) 
}


func parseGeneInput(input string) ([] string) {
    
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


func Init() {
    keggInterface = kegg.Init() 
}


func main() {

    Init() 
    

    // cmd line flags
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":8000" ,"port to run on")
    flag.Parse() 

    address := *ip + *port

    http.HandleFunc("/", indexHandler) 
    http.HandleFunc("/about", aboutHandler) 
    http.HandleFunc("/demo/", demoHandler) 

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

    fmt.Println("Webserver started on localhost ", address) 
    err := http.ListenAndServe(address, nil)
    if err != nil {
        log.Panic("Could not start webapp! ", err.Error())
    }

}
