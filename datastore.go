package main

import (
    "log"
    "flag"
    "code.google.com/p/gorest"
    "net/http"
    "nowac/kegg"
    "strings"
    "math"
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

    setScale gorest.EndPoint `method:"POST"
                             path:"/setscale/"
                             postdata:"string"`
                            

    // Dataset holding nowac data
    Dataset *Dataset

}

func (serv RestService) SetScale (PostData string) {
    log.Print("setting scale to ", PostData)

    
    serv.Dataset.setScale(PostData) 

    log.Print("---------------------------------------------")
} 

func (dataset *Dataset) setScale(scale string) {
    // changing scale to the same as before
    if(dataset.Scale == scale){
        return 
    }
    
    log.Println("--------- old scale -------- ", dataset.Scale)
    dataset.Scale = scale
    log.Println("--------- new scale -------- ",  dataset.Scale)

    
    
    var tmpIdExprs map[string][]float64
    var tmpGeneExprs map[string]map[string]*CaseCtrl

    tmpIdExprs = dataset.Exprs.IdExpression
    tmpGeneExprs = dataset.Exprs.GeneExpression

    dataset.Exprs.IdExpression = dataset.Exprs.DiffIdExpression
    dataset.Exprs.GeneExpression = dataset.Exprs.DiffGeneExpression
    

    dataset.Exprs.DiffIdExpression = tmpIdExprs
    dataset.Exprs.DiffGeneExpression = tmpGeneExprs

} 

// convert 
func log2(input []float64) [] float64 {
    new_vals := make([] float64, len(input)) 

    for i, value := range(input){
        v := math.Log2(value)
        new_vals[i] = v
    }

    return new_vals
} 

func exp2(input []float64) []float64 {

    new_vals := make([] float64, len(input)) 

    for i, value := range(input){
        v := math.Exp2(value)
        new_vals[i] = v
    }

    return new_vals
}


// Get gene expression for given gene
func (serv RestService) GeneExpression(Id string) []float64 {
    log.Print("Returning gene expression for gene ", Id)
    id := strings.Trim(Id, "hsa:")
    gene := kegg.GetGene(id)

    log.Print("hsa:",id," ==> ", gene.Name)
    
    name := strings.Split(gene.Name, ", ")[0]
    
    var ret []float64
    // return difference between case & ctrl
    for _, cc := range(serv.Dataset.Exprs.GeneExpression[name]) {
        ret = append(ret, cc.Case - cc.Ctrl) 
    } 

    return ret
}

func (serv RestService) AvgDiff(Id string) float64 {
    exprs := serv.GeneExpression(Id) 

    if(len(exprs) == 0){
        log.Print("Expression values for gene ", Id, " not found")
        return 0
    }

    avg := avg(exprs)

    log.Println("Average difference for gene ", Id, " is ",avg)
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
    
    var path = flag.String("path", "/Users/bjorn/stallo/data" , "path where data files are stored")
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":8888" ,"port to run on")

    flag.Parse()
    
    ds := NewDataset(*path) //Dataset{} // := NewDataset(*path)

    log.Print("dataset found at ", *path)
    
    ds.PrintDebugInfo()

    log.Print("Starting datastore at ", *ip, *port)
    restService := new(RestService)
    restService.Dataset = &ds

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

