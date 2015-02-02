var containers = []
var svg; 

var selected = "" 

var translates = {} 
var scales = {} 

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
                console.log("should bring",d3.select(this).attr("id"),"to front");
            }); 

    containers.push(id) 

    translates[id] = []
    scales[id] = ""

    d3.json("/pathway/"+id+"/json", function(error,graph){
        
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
                        d3.select("g#"+id).remove() 
                    }); 

            
        var node = container.append("g")
                .attr("class", "node")
                .selectAll("rect")
                .data(graph.nodes)
                .enter().append("g")
                .on("click", function(d){
                    console.log(d) 
                    // Click on pathway in vis
                    if(d.name.indexOf("path") >= 0){
                        var id = d.name.split("path:")[1];
                        // if the pathway label was clicked show info panel
                        if(d.id > 1){ 
                            pathway(id, "content", 0, 0) 
                        }
                        else { 
                            console.log("should show info panel for pathway", d.name); 
                            pathwayInfo(id); 
                        }
                    }
                }) 
                //.call(drag) 

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
                    return d.width
                })
                .attr("height", function(d){
                    return d.height;
                    })
                .style("fill", function(d){
                    if(d.name == "bg") {
                        return "url(#image)";
                    }
                    if(d.bgcolor == "none"){
                        return "#fff"
                    }
                    return d.bgcolor
                }) 
                .style("stroke", "black") 
                    

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


        }); 
    return container 

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

    var moveX = d3.event.sourceEvent.webkitMovementX;
    var moveY = d3.event.sourceEvent.webkitMovementY; 
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
    d3.select("g#"+selected).attr("transform", "translate(" + translate + ")scale(" + scale + ")");

    translates[selected] = translate
        scales[selected] = scale

}

function dragstarted(d) {
    d3.event.sourceEvent.stopPropagation();
    d3.select(this).classed("dragging", true);
}

function dragged(d) {
  //d3.select(this).attr("translate("+d3.event.translate+")"); 
          //"x", d.x = d3.event.x).attr("y", d.y = d3.event.y);
}

function dragended(d) {
    d3.select(this).classed("dragging", false);
}




