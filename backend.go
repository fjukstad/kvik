package main


import (
    "log" 
    "flag" 
    "net/http"
    "strings"
    "code.google.com/p/gorest" 
    "nowac/kegg"
) 

func main () {

    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":8080" ,"port to run on")

    flag.Parse()
    address := *ip+*port

    gorest.RegisterService(new(NOWACService)) 
    http.Handle("/", gorest.Handle()) 

    log.Println("Starting server on", address)
    err := http.ListenAndServe(address, nil) 
    if err != nil{
        log.Panic("Could not start rest-service:", err)
    }

}

type NOWACService struct {
    gorest.RestService `root:"/"
                        consumes:"application/json"
                        produces:"application/json"`

    newPathwayGraph gorest.EndPoint `method:"GET" 
                                    path:"/new/graph/pathway/{Pathways:string}"
                                    output:"string"`

}

func (serv NOWACService) NewPathwayGraph(Pathways string) string {
    addAccessControlAllowOriginHeader(serv)     
    log.Print("Pathways:", parsePathwayInput(Pathways));
    
    pws := parsePathwayInput(Pathways); 
    handlerAddress := kegg.PathwayGraphFrom(pws[0]) 

    return handlerAddress+"/"+pws[0]
    
}

func addAccessControlAllowOriginHeader (serv NOWACService) {
    // Allowing access control stuff
    rb := serv.ResponseBuilder()
    rb.AddHeader("Access-Control-Allow-Origin","*")
}

func parsePathwayInput(input string) ([] string) {
        // Remove any unwanted characters 
	a := strings.Replace(input, "%3A", ":", -1)
	a = strings.Replace(a, "&", "", -1)
	a = strings.Replace(a, "=", "", -1)
	
	// Split into separate hsa:... strings
	b := strings.Split(a, "pathwaySelect")
		
	// Clear out first empty item 
	b = b[1:len(b)]
    
    return b

}
