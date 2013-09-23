

geneIndex = 4; 
selectedGenes = window.location.href.split('/')[geneIndex];

var width = window.innerWidth,
    height = window.innerHeight - 200; 

//var width  =900,
//    height = 400; 

var svg = d3.select("body").append("svg:svg")
    .attr("width", width)
    .attr("height", height); 

var force = d3.layout.force()
    .gravity(.05)
    .distance(200)
    .charge(-100)
    .size([width, height]); 

var color = d3.scale.category10();

var url = "http://localhost:8080/api/dataset/getGenes/"+selectedGenes
var websocketUrl = "ws://localhost:3999/"

d3.json(url, function(error, json) {

    force.nodes(json.nodes) 
        .links(json.links)
        .start(); 

        var link = svg.selectAll(".link")
                .data(json.links)
                .enter().append("line")
                .attr("class", "link")
                .style("stroke-width", function(d){
                    return Math.sqrt(d.num_genes);
                })
                .style("stroke", function(d){
                    return color(d.num_genes)
                }); 

    var node = svg.selectAll(".node")
                .data(json.nodes)
                .enter().append("g")
                .attr("class", "node")
                .on("mouseover", mouseover)
                .on("mouseout", mouseout)
                .call(force.drag);


    node.append("circle") 
        .attr("r", function(d){
            if(d.size <= 0){
                return 2; 
            }
            return Math.sqrt(d.size);
        })
        .style("fill", function(d){
            console.log("COLORING:",d);
            return color(d.keggid);
        })


    node.append("text")
        .attr("dx", 23)
        .attr("dy", ".35em")
        .text(function(d) {
            console.log(d.name); 
            return d.name; 
        }); 

    force.on("tick", function() {
        link.attr("x1", function(d) { return d.source.x; })
            .attr("y1", function(d) { return d.source.y; })
            .attr("x2", function(d) { return d.target.x; })
            .attr("y2", function(d) { return d.target.y; });

    node.attr("transform", function(d) { 
         return "translate(" + d.x + "," + d.y + ")"; 
     });

});
    initWebSocket(); 
});

var socket; 
function initWebSocket(){

    socket = new WebSocket(websocketUrl);
    
    socket.onmessage = function(m) {
        console.log("Received:", m.data);
    }
    socket.error = function(m){
        console.log("WebSocketError", m.data);
    }

}


function sendMessageToServer(message) {
    message = message + "\n"
    socket.send(message)
    console.log("Sending",message,"to WebSocket @", socket.URL); 
}

function mouseover() {
  d3.select(this).select("circle").transition()
      .duration(300)
      .attr("r", 16);

  // Find id and such for selected gene
  var node = d3.select(this).select("circle").data()[0]

  //sendMessageToServer("geneId: " + node.id + " name: "+node.name)
  var jsonReply = JSON.stringify(node)
  sendMessageToServer(jsonReply)

}

function mouseout() {
  d3.select(this).select("circle").transition()
      .duration(750)
      .attr("r", function(d){
           if(d.size <= 0){
                return 2; 
            }
          return Math.sqrt(d.size);
      });

  var node = {
      id: -1,
      name: "nil",
      index: -1,
      size: -1,
      keggid: -1,
  }
  var jsonReply = JSON.stringify(node);
  sendMessageToServer(jsonReply); 


}




// Adding custom css to page 
function addCSS(cssPath) {
    linkElement = document.createElement("link");
    linkElement.rel = "stylesheet";
    linkElement.href = cssPath; 

    document.head.appendChild(linkElement);
}
addCSS("/css/pathway_graph.css"); 

