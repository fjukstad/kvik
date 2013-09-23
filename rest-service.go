package main

import (
    "code.google.com/p/gorest"
    "net/http"
    "log"
    "nowac/kegg" 
    "nowac/graph"
    "nowac/dataset"
    "encoding/json"
    "strings"
    "strconv"
    "flag"
    "time"
    "github.com/davecheney/profile"    
    "os"
)

// Global variables, I do not like this one....
var keggInterface kegg.Kegg
var keggGraph graph.Graph
var geneDataset [] dataset.Gene
var numberOfGenes = 100

func main(){

    defer profile.Start(profile.CPUProfile).Stop()

    
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":8080" ,"port to run on")

    flag.Parse()

    Init(); 
    
    gorest.RegisterService(new(NOWACService))
    http.Handle("/", gorest.Handle())
    address := *ip+*port

    log.Println("Starting server on ", address)

    err := http.ListenAndServe(address, nil)
    
    if err != nil{
        log.Panic("Error starting rest-service")
    }

    // Using SQUID CACHE
    os.Setenv("HTTP_PROXY", "http://localhost:3128")



}

type NOWACService struct{

    // Service level config
    gorest.RestService `root:"/api/" 
                        consumes:"application/json"
                        produces:"application/json"`
/*  
    geneDetails gorest.EndPoint `method:"GET" 
                                path:"/genes/{Id:int}"
                                output:"Gene"`
    
    listGenes gorest.EndPoint `method:"GET"
                                path:"/genes/"
                                output:"[]Gene"`

    listNGenes gorest.EndPoint `method:"GET"
                                path:"/ngenes/{NumGenes:int}"
                                output:"[]Gene"`
*/
    keggGraph gorest.EndPoint `method:"GET"
                               path:"/graph/"
                               output:"string"`

    geneDataset gorest.EndPoint `method:"GET"
                                path:"/genedataset/{NumGenes:int}/{NumYears:int}"
                                output:"string"`

    getGenes gorest.EndPoint `method:"GET"
                                path:"/dataset/getGenes/{Genes:string}"
                                output:"string"`
}

func Init() {
    keggInterface = kegg.Init()

    // Should possibly move this to some init thing
    keggGraph = graph.Init()
    
    t0 := time.Now() 
    // keggGraph.PopulateKeggGraph(keggInterface, numberOfGenes) 
    keggGraph.CreateKeggGraph(keggInterface, numberOfGenes) 
    t1 := time.Now() 
    log.Printf("Constructing the KEGG graph took %v\n", t1.Sub(t0))    

}

func (serv NOWACService) GeneDataset (NumGenes, NumYears int) (string) {

    addAccessControlAllowOriginHeader(serv) 
   
    var geneIds = keggGraph.GetGeneIds()
    
    
    // WARNING!!!!!!!! THIS REDUCES THE SLICE TO NUMBER OF GENES!!!!!
    geneIds = geneIds[:numberOfGenes]
    log.Print("Warning: The database contains only", numberOfGenes,"genes!") 

    log.Print(geneIds)
    
    geneDataset = dataset.DatasetFromGeneIds(geneIds, NumYears)

    b, err := json.Marshal(geneDataset) 
    if err != nil {
        log.Panic("JSON marshaling gone wrong with gene dataset") 
    }
    
    log.Print("Requets for gene dataset completed")
    return string(b) 
}

func addAccessControlAllowOriginHeader (serv NOWACService) {
    // Allowing access control stuff
    rb := serv.ResponseBuilder()
    rb.AddHeader("Access-Control-Allow-Origin","*")
}

func (serv NOWACService) KeggGraph () (string) {
    addAccessControlAllowOriginHeader(serv) 
    
    b, err := json.Marshal(keggGraph)
    if err != nil {
        log.Panic("JSON Gone wrong")
    }

    log.Println("Request for entire KEGG graph completed")
    return  string(b)
}




func (serv NOWACService) GetGenes(Genes string) (string){
    addAccessControlAllowOriginHeader(serv) 
    
    log.Println("Input genes:" , Genes) 

    genes := parseGeneInput(Genes) 

    log.Println(genes); 
    log.Println(keggGraph)

    graph := keggGraph.GetSubGraph(genes); 
    
    graph.Print()
    
    b, err := json.Marshal(keggGraph)
    if err != nil {
        log.Panic("JSON gone wrong, genes not retrieved") 
    }


    return string(b)

    
}


func keggIdsToInts(keggIds [] string) ([]int){
	var Ids [] int
	for i := range keggIds{ 
		idString := strings.Replace(keggIds[i], "hsa", "", -1)
		id, err := strconv.Atoi(idString)
		if err != nil {
			log.Println("Error in conversion")
			continue
		}
		Ids = append(Ids, id)
	}
	return Ids
}


func parseGeneInput(input string) ([] int) {
    a := strings.Replace(input, "%3A", ":", -1)
	a = strings.Replace(a, "&", "", -1)
	a = strings.Replace(a, "=", "", -1)
	
	b := strings.Split(a, "geneSelect")
	
	b = b[1:len(b)]
	
	c := keggIdsToInts(b)
	
	return c
}

/*
func (serv NOWACService) SelectGenes(Gene string)(string){
    addAccessControlAllowOriginHeader(serv) 
    log.Println("gene:", Gene)

    genes := parseGeneInput(Gene) 
    
    var parsedGenes string

    for i := range genes{
        parsedGenes = parsedGenes + genes[i]
    }

    log.Println(genes) 
    log.Println(parsedGenes) 
    
    result := "this is the reply"

    return result
}
*/

