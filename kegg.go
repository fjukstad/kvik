package kegg

import(
    //"net/http"
    "log"
    "io/ioutil"
    "strings"
    "strconv"
    "time" 
    "github.com/fjukstad/gocache"
)


type Kegg struct {
    Genes map[int] Gene
    Pathways map[int] Pathway
}


type Gene struct {
    Id int
    KEGGId string
    Name string

    Pathways map[int] Pathway
}

type Disease struct {
    KEGGId string
    Name string
}

type Drug struct {
    KEGGId string
    Name string
}

type Compound struct {
    KEGGId string
    Name string
}

type Pathway struct {
    Id int
    KEGGId string
    Name string

    // Genes [] Gene
    // Description string

    // Reference to authors should also be included
}

func Init() Kegg{
    
    genes := downloadGeneNames()
    pathways := downloadPathwayNames()
    
    createGeneGraph(genes, pathways)

    kegg := Kegg{genes, pathways} 
    
    return kegg
}

func (kegg Kegg) Test() {
    log.Println("test")

}



func CommonPathways(geneA, geneB Gene) ([] Pathway){
    
    var commonPathways []Pathway
    for _, pathway := range geneA.Pathways{
        if geneInPathway(geneB, pathway){
            commonPathways = append(commonPathways, pathway)
        }
    }
    
    return commonPathways
}

func geneInPathway(gene Gene, pathway Pathway) (bool){
    
    _, pathwayPresent := gene.Pathways[pathway.Id]
    if pathwayPresent {
        return true
    }
    return false

}

func createGeneGraph(genes map[int] Gene, pathways map[int] Pathway) {

    //url := "http://rest.kegg.jp/link/pathway/hsa"
    url := "http://localhost:8001/kegg/link-pathway-hsa"
    lines := readLinesFromURL(url) 


    for i := range lines {
        entry := strings.Split(lines[i], "\t")
        if len (entry) > 1 {
            
            keggPathwayString := strings.Split(entry[1], ":")
            keggGeneId := strings.Split(entry[0], ":")
            
            keggPathwayId := keggPathwayString[1]
            pathwayIdString := strings.Replace(keggPathwayId, "hsa","",-1)

            pathwayId, err := strconv.Atoi(pathwayIdString)
            if err != nil {
                log.Panic("strconv error: ", pathwayIdString)
            }

            geneId, err := strconv.Atoi(keggGeneId[1])
            if err != nil {
                log.Panic("String conversion gone wrong: ", keggGeneId[1])
            }

            //log.Println(lines[i])
            //log.Println("genes[",geneId,"]: ", genes[geneId], 
            //            "pathways[",pathwayId,"]:", pathways[pathwayId])

            gene := genes[geneId]
            pathway := pathways[pathwayId]
            
            gene.addPathway(pathway) 

        }
    }

}

func (gene *Gene) addPathway (pathway Pathway) {
    gene.Pathways[pathway.Id] = pathway
}


func (k Kegg) GetAllGenes() (map[int] Gene){
    return k.Genes
}

// 
func (k Kegg) GetNFirstGeneIDs(n int) ([]string){
    var genes []string 
    i := 0 

    for _, value := range k.Genes{
        genes = append(genes, value.KEGGId)     
        i += 1
        if i > n {
            break
        }
    }

    return genes
}

func readLinesFromURL(url string) (lines []string){

    response, err := gocache.Get(url)
    if err != nil {
        log.Panicf("Cannot download from url ", url)
    }

    defer response.Body.Close() // Will close response when we're done with it
    
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Panicf("error reading body")
    }

    lines = strings.Split(string(body), "\n")

    
    return lines
}

func readNLinesFromURL(numLines int, url string)(lines []string){
    
    lines = readLinesFromURL(url)
    if len(lines) > numLines{
        lines = lines[:numLines]
    }
    return lines

}

func downloadPathwayNames() (map[int] Pathway) {
    pathways := make(map[int] Pathway)
    
    //url := "http://rest.kegg.jp/list/pathway/hsa/"
    url := "http://localhost:8001/kegg/list-pathway-hsa"

    log.Print("Downloading pathway list from ", url); 

    lines := readLinesFromURL(url) 

    log.Print("lines: ", lines) 

    // Parse list of pathways containing id - name mapping
    for i:= range lines {
        pathwayString := strings.Split(lines[i], "\t")

        if len(pathwayString) > 1 {
            pathwayId := strings.Split(pathwayString[0], ":")
            
            // name := pathwayString[1]
            keggId := pathwayId[1]

            // remove hsa part of id to construct our own
            idString := strings.Replace(keggId, "hsa", "", -1) 
            id, err:= strconv.Atoi(idString)
            if err != nil {
                log.Panic("String to integer conversion cone wrong")
            }
            
            //TODO: query KEGG for a pathway

            // pathway := Pathway{id, keggId, name}
            t0 := time.Now() 

            pathway := downloadPathway(keggId) 
            
            t1 := time.Now() 
            log.Println("Downloading pathway", keggId, "took", t1.Sub(t0))
            pathways[id] = pathway
             

            log.Println("Pathway:", pathway) 
            //return pathways
        }
    }
    log.Print(" ... done") 
    return pathways
}


func downloadGeneNames() (map[int] Gene) {

    //url := "http://rest.kegg.jp/list/hsa"
    url := "http://localhost:8001/kegg/list-hsa"
    log.Println("Downloading gene names from ", url)
    lines := readLinesFromURL(url)

    genes := make(map[int] Gene)

    // Parse data from kegg 
    for i:= range lines {


        geneString := strings.Split(lines[i], "\t") 
        
        // Check for valid input
        if len(geneString) > 1 {

            geneId := strings.Split(geneString[0], ":")
                
            name := geneString[1]
            
            id, err := strconv.Atoi(geneId[1])
            if err != nil{
                log.Panicf("Integer conversion error: ",id)
            }

            keggId := "hsa"+geneId[1]
    
            pathways := make(map[int] Pathway)
            gene := Gene{id, keggId, name, pathways}
            genes[id] = gene
            
        }
    }

    log.Print(" ... done") 

    return genes
}

func downloadPathway(keggId string) (Pathway)  {
    

    p := Pathway {}

    url := "http://rest.kegg.jp/get/"+keggId


    lines := readLinesFromURL(url) 
    
    log.Print("Got pathways:")
    log.Println(lines)
    log.Print("More to come:")

    
    return p 
}

func createPathwayGraph (keggId string) (g *gographer.Graph){



}
