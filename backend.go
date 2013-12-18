package main


import (
    "log" 
    "flag" 
    "net/http"
    "strings"
    "code.google.com/p/gorest" 
    "nowac/kegg"
    "time"
    "math/rand"
    "strconv"
    "github.com/fjukstad/gocache"    
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

     getParallelVis gorest.EndPoint `method:"GET"
                            path:"/parallel/"
                            output:"string"`

    datastore gorest.EndPoint `method:"GET"
                                path:"/datastore/{...:string}"
                                output:"string"`


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
    resp, err := http.Get(URL)
    if err != nil {
        log.Print("request to datastore failed. ",err)
        serv.ResponseBuilder().SetResponseCode(404).Overide(true)
        return ":("
    }
    
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



func (serv NOWACService) GetParallelVis() string {
    addAccessControlAllowOriginHeader(serv)     
    

    code := ParallelCoordinates(10)

    log.Print("Returning parallel coordinates:", code);

    return code
}


func (serv NOWACService) GetGeneVis(Gene string) string {
    addAccessControlAllowOriginHeader(serv)     
    
    log.Print("Returning the VIS code for gene: ", Gene)

    code := Barchart() // GeneVisCode(Gene) // ParallelCoordinates(len(Gene))//GeneVisCode(Gene)
    return code
}


func (serv NOWACService) GetInfo(Items string, InfoType string) string {

    //TODO: implement different info types such as name/sequence/ etc
    
    addAccessControlAllowOriginHeader(serv)     

    if(strings.Contains(Items, "hsa")){
        // will get the first gene in the list Items. Could be more than one
        // but for starters we'll do with just one. 
        
        geneIdString := strings.Split(Items, " ")[0]
        geneId := strings.Split(geneIdString, ":")[1]

        gene := kegg.GetGene(geneId)
        return kegg.GeneToString(gene)
    }
    

    return Items;


}

func (serv NOWACService) NewPathwayGraph(Pathways string) string {
    addAccessControlAllowOriginHeader(serv)     
    
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


func Barchart() string {
    vis := `

    <style>
        .chart rect {
           fill: steelblue;
           stroke: white;
         }

    
    </style>
    <body>
    <script src="http://d3js.org/d3.v2.min.js?2.10.0"></script>
    <script
    src="https://raw.github.com/mbostock/d3/master/lib/colorbrewer/colorbrewer.js"></script>
    <script>

     function next() {
       return {
             time: ++t,
             value: v = ~~Math.max(10, Math.min(90, v + 10 * (Math.random() - .5)))
           };
         }

        var t = 1297110663, // start time (seconds since epoch)
          v = 70, // start value (subscribers)
          data = d3.range(25).map(next); // starting dataset

         var w = 20,
             h = 80;
         
        var colorScale = d3.scale.category20();
        /*
        var colorScale = d3.scale.ordinal() 
            .domain(data)
            .range(colorbrewer.YlGn[9]);
        */

         var x = d3.scale.linear()
             .domain([0, 1])
             .range([0, w]);
         
         var y = d3.scale.linear()
             .domain([0, 100])
             .rangeRound([0, h]);

         var chart = d3.select(".visman").append("svg")
            .attr("class", "chart")
            .attr("width", w * data.length - 1)
            .attr("height", h);

         chart.selectAll("rect")
             .data(data)
             .enter().append("rect")
             .attr("x", function(d, i) { return x(i) - .5; })
             .attr("y", function(d) { return h - y(d.value) - .5; })
             .attr("width", w)
             .attr("height", function(d) { return y(d.value); })
             //.style("fill", function(d, i) { return colorScale(d.value); });
             

         chart.append("line")
             .attr("x1", 0)
             .attr("x2", w * data.length)
             .attr("y1", h - .5)
             .attr("y2", h - .5)
             .style("stroke", "#ccc");

    </script>`

    return vis;


}

func GeneVisCode(gene string) string {

    
    vis := `<style>


    body {
      font: 10px sans-serif;
    }

    .bar rect {
      fill: steelblue;
      shape-rendering: crispEdges;
    }

    .bar text {
      fill: #fff;
    }

    .axis path, .axis line {
      fill: none;
      stroke: #000;
      shape-rendering: crispEdges;
    }

    </style>
    <body>
    <script src="http://d3js.org/d3.v2.min.js?2.10.0"></script>
    <script>

    // Generate an Irwin–Hall distribution of 10 random variables.
    var values = d3.range(1000).map(d3.random.irwinHall(10));

    // A formatter for counts.
    var formatCount = d3.format(",.0f");

    var margin = {top: 10, right: 30, bottom: 30, left: 30},
        width = 500 - margin.left - margin.right,
        height = 250 - margin.top - margin.bottom;

    var x = d3.scale.linear()
        .domain([0, 1])
        .range([0, width]);
    
    
    // Generate a histogram using twenty uniformly-spaced bins.
    var data = d3.layout.histogram().bins(x.ticks(20))
    (values);

    var y = d3.scale.linear()
        .domain([0, d3.max(data, function(d) { return d.y; })])
        .range([height, 0]);

    var xAxis = d3.svg.axis()
        .scale(x)
        .orient("bottom");

    var svg = d3.select(".visman").append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
      .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    var bar = svg.selectAll(".bar")
        .data(data)
      .enter().append("g")
        .attr("class", "bar")
        .attr("transform", function(d) { return "translate(" + x(d.x) + "," + y(d.y) + ")"; });

    bar.append("rect")
        .attr("x", 1)
        .attr("width", x(data[0].dx) - 1)
        .attr("height", function(d) { return height - y(d.y); });

    bar.append("text")
        .attr("dy", ".75em")
        .attr("y", 6)
        .attr("x", x(data[0].dx) / 2)
        .attr("text-anchor", "middle")
        .text(function(d) { return formatCount(d.y); });

    svg.append("g")
        .attr("class", "x axis")
        .attr("transform", "translate(0," + height + ")")
        .call(xAxis);

    </script>
    `

    return vis
    
}


func ParallelCoordinates(numGenes int) string {

    ds := GenerateDataset(150,7)

    // Header, containing all other js 
    header := `
        <div id="example" class="parcoords" style="width:450px;height:150px"></div>
        <!--- Parallel coordinates -->
        <script src="http://syntagmatic.github.io/parallel-coordinates/d3.parcoords.js"></script>
        <link rel="stylesheet" type="text/css" href="http://syntagmatic.github.io/parallel-coordinates/d3.parcoords.css">
        <script>`
    
    // dataset to be used, just random numbers now
    dataset := `var data = `+ds
    
    // rest of the vis code
    vis := `
        var pc = d3.parcoords()("#example")
          .data(data)
          .render()
          .ticks(3)
          .createAxes()
          .brushable()  // enable brushing
          .interactive()  // command line mode
        </script>
    `
    return header+dataset+vis

}

func GenerateDataset(rows, columns int) string {
    
    dataset := "[\n"
    //max := 100

    for i := 0; i < rows; i++ {
        dataset += "[" + strconv.Itoa(i) + ","
        //r := rand.New(rand.NewSource(time.Now().UnixNano()))
        /*
        for j := 0; j < columns; j++ {
            
            dataset += strconv.Itoa(r.Intn(max))
            
            if j < columns - 1 {
                dataset += ","
            }
        }
        dataset += "],\n"
        */
        
    }
    dataset += "];\n"
    
    return dataset
}


func GetGeneExpression(id int) string {

    datastore := "localhost:8888"
    
    query := "/gene/"+strconv.Itoa(id)
    url := datastore+query
    response, err := gocache.Get(url)
    
    if err != nil {
        log.Panic("could not download expression ", err)
    }

    log.Print(response)

    return "dick"

}


func GenerateRandomValues(numValues int) string {

    dataset := "["
    max := 1000
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < numValues; i++ {
        dataset += strconv.Itoa(r.Intn(max))
        if i < numValues - 1 {
            dataset += ","
        }
    }
    dataset += "];"
    return dataset
}
