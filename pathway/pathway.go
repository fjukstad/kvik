package main

import (
    "log"
    "github.com/fjukstad/gocache"
    "io/ioutil"
    "encoding/xml"
)


type KeggPathway struct {
    XMLName xml.Name            `xml:"pathway"`
    Name string                 `xml:"name,attr"`
    Org string                  `xml:"org,attr"`
    Number string               `xml:"number,attr"`
    Title string                `xml:"title,attr"`
    Image string                `xml:"image,attr"`
    Link string                 `xml:"link,attr"`
    Entries []KeggEntry         `xml:"entry"`
    Relations []KeggRelation    `xml:"relation"`
}

type KeggEntry struct {
    //  XMLName xml.Name        `xml:"entry"`
    Id string               `xml:"id,attr"`
    Name string             `xml:"name,attr"`
    Type string             `xml:"type,attr"`
    Link string             `xml:"link,attr"`
    Graphics KeggGraphics   `xml:"graphics"`

}

type KeggRelation struct {
    // XMLName xml.Name        `xml:"relation"`

    Entry1 string           `xml:"entry1,attr"`
    Entry2 string           `xml:"entry2,attr"`
    Type string             `xml:"type,attr"`
    Subtypes [] KeggSubtype `xml:"relation>subtype"`
}

type KeggGraphics struct {
    // XMLName xml.Name `xml:"graphics"`

    Name string         `xml:"name,attr"`
    Fgcolor string      `xml:"fgcolor,attr"`
    Bgcolor string      `xml:"bgcolor,attr"`
    Type string         `xml:"type,attr"`
    X string            `xml:"x,attr"`
    Y string            `xml:"y,attr"`
    Width string        `xml:"width,attr"`
    Height string       `xml:"height,attr"`

}

type KeggSubtype struct {
    // XMLName xml.Name `xml:"subtype"`
    
    Name string     `xml:"name,attr"`
    Value string    `xml:"value,attr"`
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

func (pathway *KeggPathway) Print() {
    log.Println("XMLName:", pathway.XMLName)
    log.Println("Name:", pathway.Name)
    log.Println("Org:" ,pathway.Org)
    log.Println("Title:", pathway.Title)
    log.Println("Image:", pathway.Image)
    log.Println("Link", pathway.Link)
    log.Println("Entries:" , pathway.Entries)
    log.Println("Relations:", pathway.Relations)
}

func createPathwayGraph(keggId string)  {

    url := "http://rest.kegg.jp/get/"+keggId+"/kgml"
    pw := getMap(url)
    
    pathway := new(KeggPathway)
    err := xml.Unmarshal(pw, pathway) 

    if err != nil {
        log.Panic("Could not unmarshal KGML ", err)
    }

    pathway.Print() 


    //log.Print(string(pw))
    /*
    for i := range(pathway.Relations) {
        rel := pathway.Relations[i]
        log.Println(rel)
        log.Println(rel.Entry1)
        log.Println(rel.Entry2)
    }

    for j := range(pathway.Entries) {
        ent := pathway.Entries[j]
        log.Println(ent)
    }
    */
    
}



func main() {
    createPathwayGraph("hsa05200")
}
