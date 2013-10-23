// js for retrieving vis about gene
function GetVis(gene) {

    console.log("fetching the vis code for gene ", gene);

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
    console.log("Request url:", url);
    console.log("Got response:", viscode); 

    return viscode;

}
