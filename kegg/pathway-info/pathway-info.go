package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fjukstad/kegg"
)

type Graph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	Id     string
	Name   string
	X      string
	Y      string
	Width  string
	Height string
}

type Edge struct {
	From string
	To   string
}

func main() {

	pathwayId := flag.String("id", "hsa05200",
		"The id of the pathway you're interested in")

	all := flag.Bool("all", false, "get all human pathways")
	flag.Parse()

	if *all {
		pws := kegg.GetAllHumanPathways()
		for i, pw := range pws {
			PathwayGraph(pw)
			log.Println(i, pw)
		}
	} else {
		PathwayGraph(*pathwayId)
	}
}

func PathwayGraph(pathwayId string) int {
	pathway := kegg.NewKeggPathway(pathwayId)

	nodes := make([]Node, len(pathway.Entries))
	edges := make([]Edge, len(pathway.Relations))

	for i, e := range pathway.Entries {
		info := e.Graphics
		n := Node{e.Id, info.Name, info.X, info.Y, info.Width, info.Height}
		nodes[i] = n
	}

	for j, r := range pathway.Relations {
		e := Edge{r.Entry1, r.Entry2}
		edges[j] = e
	}

	graph := Graph{nodes, edges}

	b, err := json.Marshal(graph)

	if err != nil {
		fmt.Println("Could not marshal response ", err)
		return 1
	}

	filename := pathwayId + ".json"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("could not create json file ", err)
		return 1
	}
	_, err = file.Write(b)
	if err != nil {
		fmt.Println("Could not write to json file ", err)
		return 1
	}

	file.Close()

	// Download and store pathway image
	url := pathway.Image

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Image could not be downloaded ", err)
		return 1
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read response ", err)
		return 1
	}

	filename = pathwayId + ".png"
	file, err = os.Create(filename)
	if err != nil {
		fmt.Println("Could not create image file ", err)
		return 1
	}

	_, err = file.Write(body)
	if err != nil {
		fmt.Println("Could not write image ", err)
		return 1
	}

	return 0
}
