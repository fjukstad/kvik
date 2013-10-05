package kegg

import (
    "log"
    "github.com/fjukstad/gocache"
    "io"
    "encoding/csv"
    "strings"
    "fmt"
)

type Gene struct {
    Name string
    Definition string
    Orthology string
    Organism string
    Pathway []string
    Class []string
    Position string
    DBLinks map[string]string
    Structure string
    AASEQ Sequence
    NTSEQ Sequence
}

type Sequence struct {
    length int
    Sequence string
}

func GetGene(id string) string {

    baseURL := "http://rest.kegg.jp/get/"
    url := baseURL + id
    
    response, err := gocache.Get(url)
    if err != nil{
        log.Panic("Cannot download from url:",err)
    }
    /*
    body, err := ioutil.ReadAll(response.Body)
    if err != nil{
        log.Panic("Error reading body:",err)
    }
    */

    gene := parseGeneResponse(response.Body) 
    
    log.Print(gene)

    return "hei"
}

func (g Gene) Print() {
    fmt.Println("\nGene") 
    fmt.Println("\tName:", g.Name)
    fmt.Println("\tDefinition:", g.Definition) 
    fmt.Println("\tOrthology:", g.Orthology)
    fmt.Println("\tOrganism:", g.Organism)
    fmt.Println("\tPathway:", g.Pathway)
    
    fmt.Println("")
}

func parseGeneResponse(response io.ReadCloser) Gene {
    
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
    gene := Gene{}

    for i := range records {
        
        line := strings.Split(records[i][0]," ")

        switch line[0] {
        case "NAME":
            gene.Name = strings.Join(line[8:], " ") 
        case "DEFINITION":
            gene.Definition = strings.Join(line[2:]," ") 
        case "ORTHOLOGY":
            gene.Orthology = strings.Join(line[3:], " ") 
        case "ORGANISM": 
            gene.Organism = strings.Join(line[4:], " ")
        case "PATHWAY": 

        }
    }
    gene.Print()
    return gene
}

