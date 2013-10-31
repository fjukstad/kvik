package kegg

import (
    "log"
    "github.com/fjukstad/gocache"
    "github.com/fjukstad/gographer"
    "io/ioutil"
    "encoding/xml"
    "strconv"
    "io"
    "fmt"
    "encoding/csv"
    "strings"
    "math/rand"
    "time"
)

type Pathway struct {
    Id string
    Name string
    Class string
    Pathway_Map string
    Diseases [] string
    Drugs []string
    Organism string
    Genes []string
    Compounds [] string
    
}

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
    Id string               `xml:"id,attr"`
    Name string             `xml:"name,attr"`
    Type string             `xml:"type,attr"`
    Link string             `xml:"link,attr"`
    Graphics KeggGraphics   `xml:"graphics"`

}

type KeggRelation struct {
    Entry1 string           `xml:"entry1,attr"`
    Entry2 string           `xml:"entry2,attr"`
    Type string             `xml:"type,attr"`
    Subtypes [] KeggSubtype `xml:"relation>subtype"`
}

type KeggGraphics struct {
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
    Name string     `xml:"name,attr"`
    Value string    `xml:"value,attr"`
}

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

func createPathwayGraph(keggId string, inputGraph *gographer.Graph){

    url := "http://rest.kegg.jp/get/"+keggId+"/kgml"
    // url := "http://localhost:8000/public/pathway.kgml"
    pw := getMap(url)
    
    pathway := new(KeggPathway)
    err := xml.Unmarshal(pw, pathway) 

    if err != nil {
        log.Panic("Could not unmarshal KGML ", err)
    }


    // Generate some nodes
    for j := range(pathway.Entries) {
        ent := pathway.Entries[j]
        id, _ := strconv.Atoi(ent.Id)
        name := ent.Name
        t, _ := strconv.Atoi(ent.Type) 
        size := 1
        log.Println("entry:", ent)
        inputGraph.AddNode(id,name,t,size)
    }

    // Generate some edges
    for i := range(pathway.Relations) {
        rel := pathway.Relations[i]
        source, _ := strconv.Atoi(rel.Entry1)
        target, _ := strconv.Atoi(rel.Entry2)
        weight := 19
        inputGraph.AddEdge(source, target, i, weight)  
    }
    
    return 
}

func GetPathway(id string) Pathway {
    baseURL := "http://rest.kegg.jp/get/"
    url := baseURL + id
    
    response, err := gocache.Get(url)
    if err != nil{
        log.Panic("Cannot download pathway:", err)
    }

    pathway := parsePathwayResponse(response.Body) 

    return pathway
    
}

func (p Pathway) Print() {
    fmt.Println("\nPathway") 
    fmt.Println("\tId:", p.Id)
    fmt.Println("\tName:", p.Name)
    fmt.Println("\tClass:", p.Class)
    fmt.Println("\tPathway map:", p.Pathway_Map)
    fmt.Println("\tDiseases:", p.Diseases)
    fmt.Println("\tDrugs:", p.Drugs)
    fmt.Println("\tOrganism:", p.Organism)
    fmt.Println("\tGenes:", p.Genes)
    fmt.Println("\tCompounds:", p.Compounds)
    fmt.Println("")
}

func GiveMeSomePathways() [] string {
    pw := []string{
        "hsa05200",
        "hsa04915",
        "hsa04612",
        "hsa04062",
        "hsa04660",
        "hsa04630",
        "hsa04151",
        "hsa04310",
        "hsa04662",
    }
    return pw
}

func PathwayGraphFrom(pathway string) string {
    //pw := GetPathway(pathway) 
    
    // Set up graph 
    port := getRandomPort() 
    graph := gographer.NewGraphAt(":"+strconv.Itoa(port))

    // Create a graph from pathways. Will add nodes and
    // edges and all that jazz
    createPathwayGraph(pathway, graph) 
    
    // Return address where client can retrieve graph
    // TODO: maybe the graph initialization (after wsserver setup) 
    // can be done in a go routine? 

    return graph.ServerInfo() 
}

// TODO: Should store ports used for visualization gateways

func getRandomPort() int {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return 1024 + r.Intn(45000)
}



func parsePathwayResponse(response io.ReadCloser) Pathway {


    tsv := csv.NewReader(response) 
    tsv.Comma = '\t'
    tsv.Comment = '#'
    tsv.LazyQuotes = true
    tsv.TrailingComma = true
    tsv.TrimLeadingSpace = false

    records, err := tsv.ReadAll()
    
    if err != nil {
        log.Panic("Error reading records:", err)
    }

    p := Pathway{}
    tmp := make([] string, 0) 
    current := ""
    for i := range records {
        
        line := strings.Split(records[i][0]," ")
        
        switch line[0] {
        case "ENTRY":
            a := strings.Join(line[7:], " ")
            b :=  strings.Split(a, " ")[0]
            p.Id = b
    
        case "NAME":
            p.Name = strings.Join(line[8:], " ") 
            
        case "CLASS":
            p.Class = strings.Join(line[8:], " ")

        case "PATHWAY_MAP":
            p.Pathway_Map = strings.Join(line[1:], " ")             

        case "DISEASE":
            current = "DISEASE"
        
            a := strings.Join(line[5:], " ")
            tmp = append(tmp, a)

        case "DRUG":
            p.Diseases = tmp
            tmp  = make([]string, 0)
            current = "DRUG"
            a := strings.Join(line[8:], " ")
            tmp = append(tmp, a)
            

        case "ORGANISM":
            p.Drugs = tmp

            p.Organism = strings.Join(line[4:], " ")

        case "GENE":
            current = "GENE"
            tmp = make([]string,0)
            a := strings.Join(line[8:], " ")
            b :=  strings.Split(a, " ")[0]
            tmp = append(tmp, b)

        case "COMPOUND":
            p.Genes = tmp
            current = "COMPOUND"
            tmp = make([]string, 0)
            a := strings.Join(line[8:], " ")
            tmp = append(tmp, a)

        case "REFERENCE":
            p.Compounds = tmp
            current = "REFERENCE"
            break

        default:
            if(current == "DISEASE"){
                a := strings.Join(line[12:], " ")
                tmp = append(tmp, a) 
            }
            if(current == "DRUG"){
                a := strings.Join(line[0:], " ")
                b := strings.Replace(a, "    ", "",-1)
                tmp = append(tmp, b) 
            }
            if(current == "GENE"){
                a := strings.Join(line[0:], " ")
                a = strings.Replace(a, "    ", "",-1)
                b :=  strings.Split(a, " ")[0]
                tmp = append(tmp, b)

            }
            if(current == "COMPOUND"){
                a := strings.Join(line[0:], " ")
                a = strings.Replace(a, "    ", "",-1)

                tmp = append(tmp, a)

            }

        }
    }
    return p

}


