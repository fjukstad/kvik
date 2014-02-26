package main

import (
    "nowac/kegg"
    "fmt"
    "encoding/json"
    "flag"
    "os"
)

type Graph struct { 

    Nodes [] Node
    Edges [] Edge
} 

type Node struct {
    Id string
    Name string
    X string
    Y string
    Width string
    Height string
} 

type Edge struct {
    From string
    To string
} 



func main(){

    pathwayId := flag.String("id", "hsa05200",
                "The id of the pathway you're interested in")
    flag.Parse()

    pathway := kegg.NewKeggPathway(*pathwayId) 
    
    nodes := make([]Node, len(pathway.Entries))
    edges := make([]Edge, len(pathway.Relations))

    for i, e := range pathway.Entries {
        info := e.Graphics
        n := Node{e.Id, info.Name, info.X, info.Y,info.Width, info.Height}
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
        return 
    }

    filename := *pathwayId + ".json"
    file, err := os.Create(filename)
    if err != nil {
        fmt.Println("could not create json file ", err)
        return 
    }
     _, err = file.Write(b)
    if err != nil{
        fmt.Println("Could not write to json file ", err)
    }
} 
