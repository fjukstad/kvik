
// Open pathway info panel on the left with relevant information available
function pathwayInfo(id){
    console.log("Opening info panel for pathway", id) ;

    panels.open('left');
    d3.select("#geishad-panel").style("z-index", -1).style("opacity", 0); 
    d3.select("#pathway-info-panel").style("z-index", 1).style("opacity", 100); 

    // clear out old 
    d3.select("#pathway-info-panel").html("")
    d3.select("#pathway-info-panel").on("click", function(){
        panels.close();
    });

    d3.json("/pathway/"+id+"/info", function(error,pathway){

        if(error){
            return console.warn(error);
        }

        pathway = JSON.parse(pathway) 
        
        name = pathway.Name.split("-")[0]

        panel = d3.select("#pathway-info-panel")
                  .append("div")
                  .attr("class", "panel panel-default");
    
        panel.append("div")
              .attr("class", "panel-heading")
              .append("h4")
              .text(name);
    
        panel.append("div")
              .attr("class", "panel-body")
              .text(pathway.Description) 

    }); 
    
}

// Open gene info panel on the right with relevant information available 
function geneInfo(id) {

}
