package main


import (
    "log" 
    "flag" 
    "net/http"
    "strings"
    "code.google.com/p/gorest" 
    "nowac/kegg"
//    "time"
//    "math/rand"
    "strconv"
   // "github.com/fjukstad/gocache"    
    "encoding/json"
    "io/ioutil"
    "bytes"
    "os/exec"     
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

    getInfo gorest.EndPoint `method:"GET"
                            path:"/info/{Items:string}/{InfoType:string}"
                            output:"string"`

    getGeneVis gorest.EndPoint `method:"GET"
                            path:"/vis/{Gene:string}"
                            output:"string"`


    datastore gorest.EndPoint `method:"GET"
                                path:"/datastore/{...:string}"
                                output:"string"`

    datastorePost gorest.EndPoint `method:"POST"
                                    path:"/datastore/{...:string}"
                                    postdata:"string"`
    
    pathways gorest.EndPoint    `method:"GET"
                                path:"/info/gene/{Gene:string}/pathways"
                                output:"string"`

    commonPathways gorest.EndPoint  `method:"GET"
                                    path:"/info/pathway/{Pathways:string}"
                                    output:"string"`

    pathwayGeneCount gorest.EndPoint    `method:"GET"
                                        path:"/info/gene/{Genes:string}/commonpathways"
                                        output:"string"`

    pathwayIDToName gorest.EndPoint `method:"GET"   
                                    path:"/info/pathway/{Id:string}/name"
                                    output:"string"`

    commonGenes gorest.EndPoint `method:"GET"
                                path:"/info/pathway/{Pathways:string}/commongenes"
                                output:"int"`

    resetCache gorest.EndPoint `method:"GET"
                                path:"/resetcache/"
                                output:"string"` 

}


type PWMap struct {
    Map map[string] int
}

func (serv NOWACService)  ResetCache() string {
    addAccessControlAllowOriginHeader(serv)             
    
    log.Println("!!!! CLEARING CACHE !!!!!")

    cmd := exec.Command("rm", "-rf", "cache") 
    err := cmd.Run() 
    if err != nil { 
        log.Println(err) 
    } 

    return "ok"
} 


// Returns the number of common genes shared between multiple pathways
func (serv NOWACService) CommonGenes(Pathways string) int { 
    addAccessControlAllowOriginHeader(serv)             

    pathwayList := strings.Split(Pathways, " ")
    
    allGenes := make(map[string]int) 

    // Iterate over all genes from different pathways and set their count
    for _, p := range pathwayList { 
        pw := kegg.GetPathway(p) 
        genes := pw.Genes

        for _, g := range genes { 

            count := allGenes[g]
            if count != 0 {
                allGenes[g] = count + 1
                
            } else {
                allGenes[g] = 1
            }
        }

    } 

    // From map of all genes get the ones with count larger than 1 
    var commonGenes [] string
    for k, v := range allGenes { 
        if v > 1 {
            commonGenes = append(commonGenes, k)
        }
    }

    return len(commonGenes)

} 


func (serv NOWACService) PathwayIDToName (Id string) string {
    addAccessControlAllowOriginHeader(serv)             
    return kegg.ReadablePathwayName(Id)   
}

// Returns a list of pathways and the frequency of given genes. I.e.
// how many of the given genes are represented in different pathways
// Genes is a string that looks like "hsa:123+hsa:321+..."
func (serv NOWACService) PathwayGeneCount (Genes string) string {

	PathwayMap := make(map[string] int, 0)

    geneList := strings.Split(Genes, " ")

    // for every gene get its list of pathways
    for _, g := range geneList {    
        
        geneId := strings.Split(g, ":")[1]
        gene := kegg.GetGene(geneId)
        pws := kegg.Pathways(gene)
    
        // for each of its pathways, increment the counter for number
        // of genes represented in this pathway. 
        for _, p := range pws.Pathways {
            if PathwayMap[p] != 0 {
                PathwayMap[p]++
            } else {
                PathwayMap[p] = 1
            }
        }

    }
    

    b, err := json.Marshal(PathwayMap)
    if err != nil {
        log.Panic("marshaling went bad: ",err)
    }


    return string(b)
}

func (serv NOWACService) CommonPathways(Pathways string) string {
    
    return "Not implemented yet"

}



// Will return a list of pathways for a given gene 
func (serv NOWACService) Pathways (Gene string) string {
    
    geneIdString := strings.Split(Gene, " ")[0]
    geneId := strings.Split(geneIdString, ":")[1]
    log.Println(geneId)
    gene := kegg.GetGene(geneId)
    pws := kegg.Pathways(gene)
    return kegg.PathwaysJSON(pws)

}

// Handles any requests to the Datastore. Will simply make the request to the
// datastore and return the result
func (serv NOWACService) Datastore(args ...string) string {
    
    addAccessControlAllowOriginHeader(serv)         

    requestURL := serv.Context.Request().URL.Path

    // Where the datastore is running, this would be Stallo in later versions
    datastoreBaseURL := "http://localhost:8888/"

    URL := datastoreBaseURL + strings.Trim(requestURL, "/datastore")

    // NOTE: We are not caching results here, this could have been done, but
    // since we're doing work with a test dataset caching is not done.

    //NOTE: http.GET(URL) failed when the number of these calls were really
    //frequent. now trying gocache. 
    resp, err := http.Get(URL)
    if err != nil {
        log.Print("request to datastore failed. ",err)
        serv.ResponseBuilder().SetResponseCode(404).Overide(true)
        return ":("
    }

    defer resp.Body.Close()
    
    // WARNING: int64 -> int conversion. may crash and burn if more than 2^32
    // - 1 bytes were read. Response from Datastore will typically be much
    // shorter than this, so its not an issue. 
    respLength := int(resp.ContentLength) 


    // Read the response from the body and return it as a string. 
    response := make([]byte, respLength)
    _, err = resp.Body.Read(response)
    if err != nil {
        log.Print("reading response from datastore failed. ", err)
        serv.ResponseBuilder().SetResponseCode(404).Overide(true)
        return ":("
    }

    // Set response code to what was returned from Datastore. 
    // Will ensure that if a 404 is returned by datastore this is also passed
    // along
    serv.ResponseBuilder().SetResponseCode(resp.StatusCode).Overide(false)
    
    return string(response)
}


func (serv NOWACService) DatastorePost( PostData string, varArgs ...string) {
    addAccessControlAllowOriginHeader(serv)         

    requestURL := serv.Context.Request().URL.Path

    // Where the datastore is running, this would be Stallo in later versions
    datastoreBaseURL := "http://localhost:8888/"

    URL := datastoreBaseURL + strings.Replace(requestURL, "/datastore/","",-1)
    

    postContent := bytes.NewBufferString(PostData) 

    // Perform the actual http post to the datastore
    // note that we set text as datatype. will fail miserably with anything else
    _, err := http.Post(URL, "text", postContent) 
    if err != nil { 
        log.Print("Post to datastore failed. ", err)
        serv.ResponseBuilder().SetResponseCode(500).Overide(true)
    } 
    
} 




func (serv NOWACService) GetGeneVis(Gene string) string {
    addAccessControlAllowOriginHeader(serv)     
    
    log.Print("Returning the VIS code for gene: ", Gene)

    code := GeneExpression(Gene) // Barchar() // ParallelCoordinates(len(Gene))//GeneVisCode(Gene)
    return code
}


func GeneExpression(geneid string) string {
    
    id, err := strconv.Atoi(geneid)
    if err != nil {
        log.Panic("that was not a gene id: ", geneid, " ", err)
    }
    ds := GetGeneExpression(id) 
    
    // Header, containing all other js 
    header := `
        <style>

        .chart div {
          font: 10px sans-serif;
          background-color: steelblue;
          text-align: right;
          padding: 3px;
          margin: 1px;
          color: white;
        }

        </style>
        <div class="chart"></div>
        <script src="http://d3js.org/d3.v3.min.js"></script>
        <script>`
            
    // dataset to be used, just random numbers now
    dataset := `var data = `+ds
    
    // rest of the vis code
    vis := `


    var margin = {top: 30, right: 10, bottom: 0, left: 10},
        w = $("#c1").width() - margin.left - margin.right,
        h = 170 - margin.top - margin.bottom;
        var padding = 40

        var y = d3.scale.linear()
            .domain([-d3.max(data), d3.max(data)])
            .range([0, h]);
        
        var x = d3.scale.linear()
            .domain([-d3.max(data), d3.max(data)])
            .range([0,w]);

        var xAxis = d3.svg.axis()
            .scale(x)
            .ticks(0)
            .orient("bottom");

        var yAxis = d3.svg.axis()
            .scale(y)
            .ticks(6) 
            .tickFormat(function(d) { 
                return d * -1;
            })
            .outerTickSize(0)
            .orient("left"); 
        
        var svg = d3.select(".chart")
                    .append("svg")
                    .attr("width", w)
                    .attr("height", h)
        
       
        svg.selectAll("rect")
           .data(data)
           .enter()
           .append("rect")
            .attr("x", function(d, i) {
                //console.log(d,i)
                return padding*1.2 + i * 4;  //Bar width of 20 plus 1 for padding
            })
         .attr("y", function(d) {
             if(d>0) {
                 return h - y(d) 
             }
            return h/2;  //Height minus data value
        })
        .attr("fill", function(d){
            return color(d);
        })
        
       .attr("width", 3+"px")
       .attr("height", function(d) {
            return Math.abs(y(d) - y(0));
        })
        .on("click", function(d) {
            //console.log("clicked")
            ShowBgInfo(info.Id,d)
        })

        .append("svg:title")
        .text(function(d) { 
            //console.log("hepp");
            return GetBg(info.Id, d); 
        });
        

         svg.append("g")
            .attr("class", "x axis")
            .attr("transform", "translate("+padding+","+h/2+")")
            .call(xAxis);

        
        svg.append("g")
            .attr("class", "y axis")
            .attr("transform", "translate(" + padding + ",0)")
            .call(yAxis);

    var avg =  parseFloat(AvgDiff(info.Id))
    svg.append("line")
        .attr("x1", padding)
        .attr("y1", h - y(avg))
        .attr("x2", w)
        .attr("y2", h - y(avg))
        .style("stroke", "#fab");

    var std =  parseFloat(Std(info.Id))

    var stdup = avg + std
    var stddown = avg - std

    //console.log(stdup,stddown)

    svg.append("line")
        .attr("x1", padding)
        .attr("y1", h - y(stdup))
        .attr("x2", w)
        .attr("y2", h - y(stdup))
        .style("stroke-dasharray", ("3, 3"))
        .style("stroke", "a6bbc8");

    svg.append("line")
        .attr("x1", padding)
        .attr("y1", h - y(stddown))
        .attr("x2", w)
        .attr("y2", h - y(stddown))
        .style("stroke-dasharray", ("3, 3"))
        .style("stroke", "a6bbc8");


    ypos = function(y) {
        if(y > 0){
            return h - y
        }
        return y
    } 

    var sortOrder = false;
    var sortBars = function () {
        sortOrder = !sortOrder;
        
        sortItems = function (a, b) {
            //console.log(a,b)
            if (sortOrder) {
                return a.value - b.value;
            }
            return b.value - a.value;
        };

        svg.selectAll("rect")
            .sort(sortItems)
            .transition()
            .delay(function (d, i) {
            return i ;
        })
            .duration(1000)
            .attr("x", function (d, i) {
                //console.log(d,i,x(d))
                return x(d)+4; 

        });
    } 

    d3.select("#sort").on("click", sortBars);
    
                

        </script>
    `

    return header+dataset+vis


} 


// Returns all information possible for a gene. This includes stuff
// like id,name,definition etc etc. 
func (serv NOWACService) GetInfo(Items string, InfoType string) string {

    //TODO: implement different info types such as name/sequence/ etc
    
    addAccessControlAllowOriginHeader(serv)     

    if(strings.Contains(Items, "hsa")){
        // will get the first gene in the list Items. Could be more than one
        // but for starters we'll do with just one. 
        
        geneIdString := strings.Split(Items, " ")[0]
        geneId := strings.Split(geneIdString, ":")[1]

        gene := kegg.GetGene(geneId)

        //gene.Pathways = kegg.ReadablePathwayNames(gene.Pathways) 

        return kegg.GeneJSON(gene)
    }
    

    return Items;


}

func (serv NOWACService) NewPathwayGraph(Pathways string) string {
    addAccessControlAllowOriginHeader(serv)     
    
    pws := parsePathwayInput(Pathways); 
    log.Print(Pathways)
    handlerAddress := kegg.PathwayGraphFrom(pws[0]) 

    return handlerAddress+"/"+pws[0]
    
}

func addAccessControlAllowOriginHeader (serv NOWACService) {
    // Allowing access control stuff
    rb := serv.ResponseBuilder()
    if serv.Context != nil {
        rb.AddHeader("Access-Control-Allow-Origin","*")
    }
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


func GetGeneExpression(id int) string {

    datastore := "http://localhost:8888"
    
    query := "/gene/"+strconv.Itoa(id)
    url := datastore+query
    response, err := http.Get(url)
    
    if err != nil {
        log.Panic("could not download expression ", err)
    }

    defer response.Body.Close()

	exprs, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Panic("Could not read expression ", err) 
    } 
    

    
    return string(exprs)

}



