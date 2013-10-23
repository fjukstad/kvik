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

    getInfo gorest.EndPoint `method:"GET"
                            path:"/info/{Items:string}/{InfoType:string}"
                            output:"string"`

    getGeneVis gorest.EndPoint `method:"GET"
                            path:"/vis/{Gene:string}"
                            output:"string"`
                      

}

func (serv NOWACService) GetGeneVis(Gene string) string {
    addAccessControlAllowOriginHeader(serv)     
    
    log.Print("Returning the VIS code for gene: ", Gene)
    code := GeneVisCode(Gene)
    log.Print("Returning:", code) 
    return code
}


func (serv NOWACService) GetInfo(Items string, InfoType string) string {

    //TODO: implement different info types such as name/sequence/ etc
    
    addAccessControlAllowOriginHeader(serv)     

    log.Println("now fetchign items", Items);
    log.Println("for info type:", InfoType);

    if(strings.Contains(Items, "hsa")){
        log.Println("this here is a gene!");
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
    var data = d3.layout.histogram()
        .bins(x.ticks(20))
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
