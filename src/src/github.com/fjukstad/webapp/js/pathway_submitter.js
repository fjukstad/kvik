

$(function() {
    $('#pathwayFormSubmit').click(function(e){
        e.preventDefault();
        $.get(getPathways(),function(data) {
              $('.result').html(data);
            });
        }); 
});



function getPathways() {
    var pathways = $('#pathwaySelect').serialize()
    var result = "/browser/"+pathways
    window.location.assign(result)
    
   // window.open(result,'name','toolbar=0,status=0,height=700,width=1300');
    
    return result; 
}

function ReadableInput() {
    var sel = document.getElementById('pathwaySelect');
    var opts = sel.options;
    for(var opt, j = 0; opt = opts[j]; j++) {
        opt.innerHTML = GetPathwayName(opt.innerHTML)
    }
} 


function GetPathwayName(id) {
    var baseURL = "http://"+window.location.hostname+":8080/info/pathway/"
        url =  baseURL+id+"/name"

    var name
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            name = data;
        }
    });
    return name;

}


$(function() {
    $('#geneFormSubmit').click(function(e){
        e.preventDefault();
        $.get(getGenes(),function(data) {
              $('.result').html(data);
            });
        }); 
});



function getGenes() {
    var genes = $('#geneSelect').serialize()
    var result = "/browser/"+genes
    window.location.assign(result)
     return result; 
}

var pcolor = d3.scale.linear()
    .domain([0.0083, 0.00915, 0.010])
    .range(colorbrewer.RdBu[3]);
    //.range(["yellow", "white", "blue"]);

// When window load fetch names for the pathways in the list.
window.onload = function() {
    ReadableInput() 
    var rawgenes = GetGenes().Genes

    console.log(rawgenes) 

    var genelist = "" 
    var genes = [] 
    for(var i = 0; i < rawgenes.length; i++){
        var gene = rawgenes[i].replace(/\"/g, "") // removing "" 
        if (i > 1) { 
            genes.push(gene) 
            genelist = genelist + "+" + gene
        } else {
            genes.push(gene)
            genelist = gene 
        } 
    }


    var dataset = GetPValues(genelist).Result
    var d_keys = Object.keys(dataset) 

    console.log(dataset) 
    
    var tr = d3
        .select("table#geneselect")
        .selectAll("tr")
        .data(d_keys)
        .enter().append("tr")
        .style("line-height", "15px") 
        
        
        tr.append("td")
        .attr("class", "col-md-1") 
        .html(function(d){
            
            return '<a style="font-size:15px" href="http://'+window.location.hostname+'/browser/geneSelect='+d+'">'+d+"</a>";
        }); 

    tr.append("td") 
        .attr("class", "col-md-1") 
        .append("svg")
                .attr("width", 100)
                .attr("height", 10)
                .style("float", "left")
        .append("svg:a")
        .attr("xlink:href",function(d){
            return "http://"+window.location.hostname+"/browser/geneSelect="+d;
        }) 
        .append("rect") 
        .attr("width", function(d){
            if (dataset[d] == "NA"){
                return 0
            }
            return dataset[d] * 10000;
        })
        .attr("height", 10) 
        .style("fill", function(d){
            console.log(parseFloat(dataset[d]))
            // significant
            if(parseFloat(dataset[d]) < 0.009978){
                return "#67a9cf"

            } else {
                return "#ef8a62" 
            }
        })

    var h = $("select#pathwaySelect").height() + 18
    $("#geneselect").height(h); 

    /*
    var h = $("select#geneSelect").height(),
        w = $("select#geneSelect").width(); 

    var formh = $("#genediv").height() 
    var offset = 66 + "px" 
    console.log(h,w,formh)

    var div = d3.select("#genediv").append("div")
        .attr("id", "outer")
        .style("border", "1px solid black")
        .style("height", h)
        .style("width", w)
        .style("margin", 0)
        .style("padding", 0) 
        .style("z-index","-1") 
        .style("top", offset)
        .style("position", "absolute")


    var div2 = div.append("div")
        .style("width", w)
        .style("height", "100")
                .style("z-index","-1") 

        .style("overflow", "y-scroll");
;


    var dataset = [10,20,30,40,50,60,70,80,90];

    var svg = div2.append("svg")
                .attr("width", w)
                .attr("height", "1000")
                        .style("z-index","-1"); 

    var rect = svg.selectAll("rect") 
        .data(dataset)
        .enter()
        .append("rect");
    
    rect.attr("x", 250)
        .attr("y", function(d){
                return d *5;
                })
        .attr("width", function(d){
                return d;
                })
        .attr("height", 10) 
        .style("fill", "black") 
    
    */

}



