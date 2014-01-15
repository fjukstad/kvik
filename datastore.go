package main

import (
    "log"
    "flag"
    "code.google.com/p/gorest"
    "net/http"
    "nowac/kegg"
    "strings"
/*
    "runtime/pprof"
    "os"
*/
)

type RestService struct  {

    // REST service details
    gorest.RestService `root:"/"
                        consumes:"application/json"
                        produces: "application/json" `
    
    geneExpression gorest.EndPoint `method:"GET"
                                    path:"/gene/{Id:string}"
                                    output:"[]float64"`

    avgDiff gorest.EndPoint `method:"GET"
                            path:"/gene/{Id:string}/avg"
                            output:"float64"`

    // Dataset holding nowac data
    Dataset Dataset

}

// Get gene expression for given gene
func (serv RestService) GeneExpression(Id string) []float64 {
    log.Print("Returning gene expression for gene ", Id)
    id := strings.Trim(Id, "hsa:")
    gene := kegg.GetGene(id)
    log.Print("hsa:",id," ==> ", gene.Name)
    name := strings.Split(gene.Name, ", ")[0]
    
    exprs := serv.Dataset.Exprs.GeneExpression[name]
    return exprs 
}

func (serv RestService) AvgDiff(Id string) float64 {

    //log.Print("Average difference in gene expression for gene ", Id[0] )


    //log.Print("Expression: ", serv.Dataset.Exprs.IdExpression[Id[0]])
    //log.Print(serv.Dataset.Exprs.Genes)

    id := strings.Trim(Id, "hsa:")
    gene := kegg.GetGene(id)

    log.Print("hsa:",id," ==> ", gene.Name)
    
    name := strings.Split(gene.Name, ", ")[0]

    log.Print(name) 

    

    exprs := serv.Dataset.Exprs.GeneExpression[name]
        
    if(len(exprs) == 0){
        log.Print("Expression values for gene ", name, " not found")
        return 0
    }

    avg := avg(exprs)

    log.Print("Avg expression for gene ", name," = ",avg)
    return avg
}

func avg(nums [] float64) float64 {
    
    var total float64

    for _, num := range(nums) {
        total += num
    }

    return total / float64(len(nums))

} 

func main() {
    
    var path = flag.String("path", "data" , "path where data files are stored")
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":8888" ,"port to run on")

    flag.Parse()
    
    ds := NewDataset(*path) //Dataset{} // := NewDataset(*path)

    log.Print("dataset found at ", *path)
    
    ds.PrintDebugInfo()

    log.Print("Starting datastore at ", *ip, *port)
    restService := new(RestService)
    restService.Dataset = ds

    gorest.RegisterService(restService)

    http.Handle("/", gorest.Handle())
    http.ListenAndServe(*port, nil)



/*
    
// Profiling
// $ go tool pprof /home/bfj001/master/src/bin/datastore memprofile.prof 
// (pprof) top5
    f, err := os.Create("memprofile.prof")
    if err != nil {
        log.Fatal(err)
    }
    pprof.WriteHeapProfile(f)
    f.Close()
    return
*/

}

