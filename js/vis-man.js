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
        cache: true,
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
