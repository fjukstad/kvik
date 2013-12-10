package main

import (
    "log"
    "flag"
    "code.google.com/p/gorest"
    "net/http"
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
                                    path:"/gene/{Id:int}"
                                    output:"[]float64"`


    // Dataset holding nowac data
    Dataset Dataset

}

// Get gene expression for given gene
func (serv RestService) GeneExpression(Id int) []float64 {
    log.Print("Returning gene expression for gene ", Id)
    output := serv.Dataset.Exprs.GeneExpression[serv.Dataset.Exprs.Genes[0]]
    return output 

}

func main() {
    
    var path = flag.String("path", "data" , "path where data files are stored")
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":8888" ,"port to run on")

    flag.Parse()
    
    ds := NewDataset(*path)
    
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

