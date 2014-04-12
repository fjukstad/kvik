package main

import (
    "log"
    "flag"
    "code.google.com/p/gorest"
    "net/http"
    "nowac/kegg"
    "strings"
    "math"
    "runtime/pprof"
    "os"
    "github.com/fjukstad/rpcman" // rpc to the statistics man
    "strconv" 
    "encoding/json" 
    "io/ioutil"
    
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

    std gorest.EndPoint `method:"GET"
                            path:"/gene/{Id:string}/stddev"
                            output:"float64"`

    variance gorest.EndPoint `method:"GET"
                            path:"/gene/{Id:string}/vari"
                            output:"float64"`

    setScale gorest.EndPoint `method:"POST"
                             path:"/setscale/"
                             postdata:"string"`
    
    bg gorest.EndPoint `method:"GET"
                        path:"/gene/{GeneId:string}/{Exprs:string}/bg"
                        output:"string"`

    // Dataset holding nowac data
    Dataset *Dataset

    // RPC Server for performing statistics 
    RPC rpcman.RPCMan
    
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

    if(gene.Name == "") { 
        log.Println("Gene with id ",Id," not found in database")

        // return slice with all zeros
        emptySlice := make([]float64, len(serv.Dataset.Exprs.Genes))
        return emptySlice
    }

    log.Println("hepp")
    
    name := strings.Split(gene.Name, ", ")[0]
    
    var ret []float64
    // return difference between case & ctrl
    for _, cc := range(serv.Dataset.Exprs.GeneExpression[name]) {
        ret = append(ret, cc.Case - cc.Ctrl) 
    } 

    return ret
}

// Get standard deviation for expression values of a given gene
func (serv RestService) Std(GeneId string) float64 { 
    exprs := serv.GeneExpression(GeneId) 

    if(len(exprs) == 0){
        log.Print("Expression values for gene ", GeneId, " not found")
        return 0
    }
    ret, err := serv.RPC.Call("std",exprs) 
    if err != nil { 
        log.Println("RPC FAILED", err) 
        return 0
    } 

    std, ok := ret.(float64) 
    if !ok{ 
        log.Println("conversion to float64 went bad: ",ret) 
        return 0
    } 

    log.Println("Standard deviation for expression of gene ", GeneId, " is ",std)
    return std


} 

// Get variation for expression values of a given gene
func (serv RestService) Variance(GeneId string) float64{ 
    exprs := serv.GeneExpression(GeneId) 

    if(len(exprs) == 0){
        log.Print("Expression values for gene ", GeneId, " not found")
        return 0
    }

    ret, err := serv.RPC.Call("var",exprs) 
    if err != nil { 
        log.Println("RPC FAILED", err) 
        return 0
    } 

    variance, ok := ret.(float64) 
    if !ok{ 
        log.Println("conversion to float64 went bad: ",ret) 
        return 0
    } 

    log.Println("Variance for expression of gene ", GeneId, " is ",variance)
    return variance

} 

func (serv RestService) AvgDiff(Id string) float64 {
    exprs := serv.GeneExpression(Id) 

    if(len(exprs) == 0){
        log.Print("Expression values for gene ", Id, " not found")
        return 0
    }

    // avg := avg(exprs)
    
    ret, err := serv.RPC.Call("mean",exprs) 
    if err != nil { 
        log.Println("RPC FAILED", err) 
    } 
    avg, ok := ret.(float64) 
    if !ok{ 
        log.Println("conversion to float64 went bad: ",ret) 
    } 


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


// Find dataset id that has the given expression value
func expressionToId(dataset *Dataset, GeneId, Exprs string) string { 
    
    id := strings.Trim(GeneId, "hsa:")
    gene := kegg.GetGene(id)

    if(gene.Name == "") { 
        log.Println("Gene with id ",GeneId," not found in database")
        return "" 
    }
    
    name := strings.Split(gene.Name, ", ")[0]
    
    dsId := ""

    exprsVal, err := strconv.ParseFloat(Exprs,64) 
    if err != nil { 
        log.Println("could not convert ",Exprs,"to float")
        return ""
    } 

    // return difference between case & ctrl
    for i, cc := range(dataset.Exprs.GeneExpression[name]) {
        ex := cc.Case - cc.Ctrl
        if(ex == exprsVal) { 
            dsId = i 
            break
        } 

    } 

    return dsId

} 

func (serv RestService) Bg(GeneId, Exprs string) string { 
    dsId := expressionToId(serv.Dataset, GeneId, Exprs) 

    bg := serv.Dataset.Bg.IdInfo[dsId]


    b, err := json.Marshal(bg)
    if err != nil {
        log.Print("marshaling went bad: ",err)
        return ""
    }

    return string(b)
} 

func Init(path string) *RestService { 
    ds := NewDataset(path) //Dataset{} // := NewDataset(*path)

    log.Print("dataset found at ", path)
    
    ds.PrintDebugInfo()

    restService := new(RestService)
    restService.Dataset = &ds

    // connect to statistics engine that will run statistics and that 
    rpcaddr := "tcp://localhost:5555" // "ipc:///tmp/datastore/0" //"tcp://localhost:5555"
    restService.RPC = rpcman.Init(rpcaddr)
    restService.RPC.Connect() 
    
    return restService
} 

func main() {

    // enable debugging
    
        log.SetOutput(ioutil.Discard) 
    //defer log.SetOutput(os.Stdout) 
    
    var path = flag.String("path", "/Users/bjorn/stallo/data" , "path where data files are stored")
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":8888" ,"port to run on")

    flag.Parse()
    
    restService := Init(*path)
    log.Print("Starting datastore at ", *ip, *port)
    
    f, err := os.Create("memprofile.prof")
    if err != nil {
        log.Print(err)
        return
    }
    pprof.WriteHeapProfile(f)
    f.Close()

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


