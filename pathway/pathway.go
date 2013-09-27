package main

import (
    "log"
   // "github.com/fjukstad/gographer"
    "github.com/fjukstad/gocache"
    "io/ioutil"
    "encoding/xml"
)


type KeggPathway struct {
    XMLName xml.Name `xml:"pathway"`
    Name string    `xml:"name,attr"`
    Org string      `xml:"org,attr"`
    Number string   `xml:"number,attr"`
    Title string    `xml:"title,attr"`
    Image string    `xml:"image,attr"`
    Link string     `xml:"link,attr"`
    Entries []KeggEntry `xml:"entry"`
    Relations []KeggRelation `xml:"relation"`
}

type KeggEntry struct {
    XMLName xml.Name `xml:"entry"`
    Id string
    Name string
    Type string
    Link string
    Graphics KeggGraphics

}

type KeggRelation struct {
    XMLName xml.Name `xml:"relation"`

    Entry1 string `xml:"entry1"`
    Entry2 string
    Type string
    Subtypes [] KeggSubtype
}

type KeggGraphics struct {
    XMLName xml.Name `xml:"graphics"`

    Name string
    fgcolor string
    bgcolor string
    Type string
    X string
    Y string
    Width string
    Height string

}

type KeggSubtype struct {
    XMLName xml.Name `xml:"subtype"`
    
    Name string
    Value string
}
/*
func downloadPathway(keggId string) (KeggPathway) {
   
    p := Pathway {}

    url := "http://rest.kegg.jp/get/"+keggId

    lines := readLinesFromURL(url) 
    
    log.Print("Got pathways:")
    log.Println(lines)
    log.Print("More to come:")

    return p  
}
*/

func getMap(url string) ([]byte) {
    response, err := gocache.Get(url)
    if err != nil {
        log.Panic("Could not download pathway kgml:",err)
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Panic("KGML body could not be read:", err) 
    }
    return body
}

func createPathwayGraph(keggId string)  {

    url := "http://rest.kegg.jp/get/"+keggId+"/kgml"
    pw := getMap(url)
    
    pathway := new(KeggPathway)
    err := xml.Unmarshal(pw, pathway) 

    if err != nil {
        log.Panic("Could not unmarshal KGML ", err)
    }

    log.Println(pathway.Name)
/*
    //log.Print(string(pw))
    for i := range(pathway.Relations) {
        rel := pathway.Relations[i]
        log.Println(rel.Entry1)
        log.Println(rel.Entry2)
    }
*/
    
}



func main() {
    createPathwayGraph("hsa05200")
}
