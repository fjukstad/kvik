function geneBars(element, data){
    console.log("HEHEHE") 
    var genePanel = newPanelWithHeader("#gene-info-panel", "Gene Expression") 

    var w = 250
    var h = 250
    var padding = 20;

    var y = d3.scale.linear()
        .domain([-1, 1])
        .range([0, h]);

    var x = d3.scale.linear()
        .domain([0, data.length])
        .range([0,w]);

    var svg = genePanel.append("svg")
                    .attr("width", w)
                    .attr("height", h+30)
                    .style("padding-left", padding-5)
                    .style("magin-top", padding)
                    .style("padding-top", padding/3)
                    .style("padding-bottom", padding/3)
    
    var yAxis = d3.svg.axis()
        .scale(y)
        .ticks(6) 
        .tickFormat(function(d) { 
            return d * -1;
        })
        .outerTickSize(0)
        .orient("left"); 

    var xAxis = d3.svg.axis()
        .scale(x)
        .ticks(0)
        .orient("bottom");
    
    svg.selectAll("rect")
        .data(data) 
        .enter()
        .append("rect")
        .attr("x", function(d,i){
            return padding + i * 4
        })
        .attr("y", function(d){
             if(d>0) {
                return h - y(d) 
             }
             return h/2; 
        })
        .attr("fill", function(d){
            console.log(d, color(d), color(parseFloat(d))) 
            return color(d)
        }) 
        .attr("width", 3)
        .attr("height", function(d){
            return Math.abs(y(d) - y(0));
        })
        .append("svg:title")
        .text(function(d) { return d});

        svg.append("g")
            .attr("class", "y axis")
            .attr("fill", "black") 
            .attr("transform", "translate(" + padding + ",0)")
            .call(yAxis);

         svg.append("g")
            .attr("class", "x axis")
            .attr("transform", "translate("+padding+","+h/2+")")
            .call(xAxis);
} 
