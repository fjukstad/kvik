var containers = []
var svg; 

var selected = "" 

var translates = {} 
var scales = {} 

var visibleNodes = []
var color
var nodeId = 0; 
    
function pathway(id, element, h, w){ 

    // First check if the pathway is already there
   try {
       d3.select("g#"+id).attr("id");
       swal({
            title: "Pathway already in view!"
        }) 
       return
   } 
   catch(err){
   } 
    color = d3.scale.linear()
        .domain([0.8351978, 0, -0.8310282])
        //.domain([1.0, 0.5, 0.0, -0.5, -1.0])
        //.range(colorbrewer.RdYlBu[3]);
        .range(["red","white", "green"])

    var margin = {top: 5, right: 5, bottom: 5, left: 5},
        width = w - margin.left - margin.right,
        height = h - margin.top - margin.bottom; 
    
    var zoom = d3.behavior.zoom()
        .scaleExtent([-10, 10])
        .on("zoom."+id, zoomed);

    var drag = d3.behavior.drag()
        .origin(function(d) {
            return d;
        })
        .on("dragstart", dragstarted)
        .on("drag", dragged) 
        .on("dragend", dragended);

    if(typeof svg === 'undefined') { 
        svg = d3.select(element).append("svg") 
                    .attr("width", width+margin.left+margin.right)
                    .attr("height", height+margin.top+margin.bottom)
                    .append("g")
                    .attr("transform", "translate(" + margin.left + "," + margin.right + ")")
                    .call(zoom)
    }

    var container = svg.append("g").attr("class", "pathway" ).attr("id", id); 

    container.on("mouseover", function (d) {
                selected = d3.select(this).attr("id"); 
            })
            .on("mouseout", function() { 
                selected = "" 
            }) 
            .on("click", function(d) {
            }); 

    containers.push(id) 

    translates[id] = []
    scales[id] = ""

    d3.json("/pathway/"+id, function(error,graph){
        
        if(error){
            return console.warn(error);
        }
        
        var bg = graph.nodes[0]
        var image = container.append("image")
                        .attr("height", bg.height)
                        .attr("width", bg.width)
                        .attr("xlink\:href", function(d){
                                return bg.description
                        }); 
    
        container.append("circle")
                .attr("cx", function(){
                    return bg.width + 5;
                })
                .attr("cy", -2.5) 
                .attr("r", 5)
                .style("fill", "red") 
                .on("click", function(){
                        d3.select("g#"+id).remove(); 
                    }); 
                

        var node = container.append("g")
                .attr("class", "node")
                .selectAll("rect")
                .data(graph.nodes)
                .enter().append("g")
                .on("click", function(d){
                    // Click on pathway in vis
                    if(d.name.indexOf("path") >= 0){
                        var id = d.name.split("path:")[1];
                        if(d.id < 1 || d.y == 58){ 
                            // TODO 
                        }
                        else { 
                            pathway(id, "content", 0, 0) 
                        }
                        return
                    }
                    if(d.shape == "rectangle"){
                        var id = d.name.split(" ")[0]
                        if(id === "bg"){
                            return
                        } else { 
                            var url = "#" ;
                            //window.open(url, '_blank');
                        }
                    }
                }) 

        node.append("rect")
                .attr("x", function(d){
                    if(d.name == "bg"){
                        return 0
                    }
                    return d.x - d.width/2;
                })
                .attr("y", function(d){
                    if(d.name == "bg"){
                        return 0
                    }
                    return d.y - d.height/2;
                })
                .attr("width", function(d){
                    if(d.name == "undefined"){
                        return 0
                    }                     
                    return d.width
                })
                .attr("height", function(d){
                    if(d.name == "undefined"){
                        return 0
                    }                     
                    return d.height;
                })
                .attr("class", function(d){
                    if(d.name.indexOf("hsa:") >=0){
                        return "gene"
                    }
                    return "" 
                })
                .style("fill", function(d){
                    if(d.name == "bg") {
                        return "url(#image)";
                    }
                    if(d.bgcolor == "none"){
                        return "#fff"
                    }
                    if(d.bgcolor == "#FFFFFF" || d.fc == ""){
                        return "#c9c9c9"
                    }
                    return color(parseFloat(d.fc))
                }) 
                
                .style("stroke", "black")
                
                .attr("id", function(d){
                    var id = d.name.split(" ")[0]
                    id = id.replace(":","")
                    return id
                }) 
                .attr("nodeid", function(d){
                    d.nodeId = nodeId; 
                    nodeId = nodeId + 1; 
                    
                    return d.nodeId;
                }) 
                    
            
        node.append("text")
            .attr("x", function(d){
                 return d.x-d.width/2.5;
             })
            .attr("y", function(d){
                 return d.y + d.height/3;
             })
            .text(function(d){
                if(d.shape == "circle" || d.name =="bg"){
                    return ""
                }
                return d.description;
            }) 

        highlightGene(oldgene) 

        }); 
    
        
    return container 

}

var oldgene = "" 
function highlightGene(id){
    
    if(id.indexOf(":") > -1){
        id = id.replace(":","")
    }

    try { 
        d3.selectAll("rect#"+oldgene)
            .attr("width", function(d){
                return d.width;
            })
            .attr("heigth", function(d){
                return d.height;
            })
          .style("stroke", "black")
          .style("stroke-width", 1); 
    } catch(err) {
    }

    try { 
        var stroke = 5; 
        d3.selectAll("rect#"+id)
            .attr("width", function(d) {
                return d.width + stroke;
            })
            .attr("height", function(d){
                return d.height + stroke;
            })
            .style("stroke", "#e7298a")
            .style("stroke-width", stroke); 
    } catch(err){
    }
    
    oldgene = id; 
}



function zoomed() {
    
    translate = translates[selected]
    scale = scales[selected] 
    
    if(scale == ""){
        translates[selected] = [0,0]
        scales[selected] = 1;

        translate = translates[selected]
        scale = scales[selected] 
    }

    var moveX = d3.event.sourceEvent.movementX;
    var moveY = d3.event.sourceEvent.movementY; 
    var clientX = d3.event.sourceEvent.clientX;
    var clientY = d3.event.sourceEvent.clientY;
    var scroll = d3.event.sourceEvent.wheelDelta; 



    if(typeof scroll !== 'undefined'){ 
        var newscale =  scale + (scroll * 0.001); 
        if(newscale > 0.05){
            scale = newscale
        }
        // TODO: Do not move when scaling
    } else {  
        if(typeof translate === 'undefined'){
            return;
        } else { 
            translate[0] = translate[0] + moveX;
            translate[1] = translate[1] + moveY;
        }
    }
    try {
        d3.select("g#"+selected).attr("transform", "translate(" + translate + ")scale(" + scale + ")");
    } catch(err){
    }

    translates[selected] = translate
        scales[selected] = scale

}

function dragstarted(d) {
    d3.event.sourceEvent.stopPropagation();
    d3.select(this).classed("dragging", true);
}

function dragged(d) {
}

function dragended(d) {
    d3.select(this).classed("dragging", false);
}


