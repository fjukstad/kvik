// js for retrieving vis about gene
function GetVis(gene) {


    var baseURL = "http://"+window.location.hostname+":8080/vis/";

    url = baseURL + gene

    var viscode
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            viscode = data;
        }
    });

    return viscode;

}

// Will download a parallel coordinate plot 
function GetParallelVis() {


    var baseURL = "http://"+window.location.hostname+":8080/parallel";

    url = baseURL

    var viscode
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            viscode = data;
        }
    });

    return viscode;

}

// Fetches avg diff from datastore
function AvgDiff(gene) {
    
    var baseURL = "http://"+window.location.hostname+":8080/datastore/gene/";
    url = baseURL + gene + "/avg"

    var avg
    $.ajax({
        async: false,
        cache: false,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            avg = data;
        }
    });

    console.log("got avg diff: ", avg)

    return avg;

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

    console.log("name: ", name)

    return name;

}



function GetCommonGenes(ids) {
    var baseURL = "http://"+window.location.hostname+":8080/info/pathway/"
        url =  baseURL+ids+"/commongenes"

    var num
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
             num = data;
        }
    });

    console.log("Number of common genes: ", num)

    return num;

}

function updateColor(scale) { 
    
    if(scale == "log") { 
        console.log("log scale") 
        color = d3.scale.linear()
            .domain([0, 1, 10])
            .range(["white", "yellow", "purple"]);
    } 
    else { 
        console.log("abs scale") 
        color = d3.scale.linear()
            .domain([0, 1, 1500])
            .range(["white", "yellow", "purple"]);
    } 



} 
