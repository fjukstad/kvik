
// Open pathway info panel on the left with relevant information available
function pathwayInfo(id){
    
    initPanels('left', "#geishad-panel", "#pathway-info-panel"); 
    
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

    initPanels('right', "#search-panel", "#gene-info-panel");

    d3.json("/gene/"+id+"/info", function(error,gene){
        console.log(gene,id)
        gene = JSON.parse(gene) 

        console.log(gene) 
        panel = d3.select("#gene-info-panel")
                  .append("div")
                  .attr("class", "panel panel-default");
        
        panel.append("div")
              .attr("class", "panel-heading")
              .append("h4")
              .text(gene.Name);
    
        panel.append("div")
              .attr("class", "panel-body")
              .text(gene.Definition) 


    }); 
}

// side: left or right
// hidden: which panel to hide
// focus: which panel to put in focus
function initPanels(side, hidden, focus) {

    panels.open(side);
    d3.select(hidden).style("z-index", -1).style("opacity", 0); 
    d3.select(focus).style("z-index", 1).style("opacity", 100); 

    // clear out old 
    d3.select(focus).html("")

    // close panel on click 
    d3.select(focus).on("click", function(){
        panels.close();
    });

} 
