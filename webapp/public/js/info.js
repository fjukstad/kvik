
// Open pathway info panel on the left with relevant information available
function pathwayInfo(id){
    
    initPanels('left', "#geishad-panel", "#pathway-info-panel"); 
    
    d3.json("/pathway/"+id+"/info", function(error,pathway){

        if(error){
            return console.warn(error);
        }

        pathway = JSON.parse(pathway) 
        
        name = pathway.Name.split("-")[0]

        addPanel("#pathway-info-panel", name, pathway.Description)
        
    }); 
    
}

// Open gene info panel on the right with relevant information available 
function geneInfo(id) {

    initPanels('right', "#search-panel", "#gene-info-panel");

    d3.json("/gene/"+id+"/info", function(error,gene){
        console.log(gene,id)
        gene = JSON.parse(gene) 

        console.log(gene) 
        panel = addPanel("#gene-info-panel", gene.Name, gene.Definition) 
        
        var pathways = newPanelWithHeader("#gene-info-panel", "Pathways") 
        console.log(gene.Pathways) 
        
        var list = pathways.append("div")
                .attr("class", "panel-body")
                .append("ul")
                .selectAll("li")
                .data(gene.Pathways)
                .enter()
                .append("li")
                .append("a")
                .attr("id", function(d){
                    return d;
                })
                .text(function(d){
                    var resp = $.get("/pathway/"+d+"/name", function(data){
                        var name = data.split(" - Homo")[0];
                        d3.select("a#"+d).text(name); 
                        return data;
                    })
                })
                .on("click", function(d){
                    pathway(d, "content", 0,0);

                }); 

    }); 
}

function addPanel(element, header, body) {

    var panel = newPanelWithHeader(element, header) 
        
    panel.append("div")
          .attr("class", "panel-body")
          .text(body) 
          
    return panel
}

function newPanelWithHeader(element, header) { 
    
    var panel = d3.select(element)
              .append("div")
              .attr("class", "panel panel-default");
    
    panel.append("div")
          .attr("class", "panel-heading")
          .append("h4")
          .text(header);

    return panel
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
