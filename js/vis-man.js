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

    //console.log(gene,"got avg diff: ", avg)

    return avg;

}

function AvgDiffs(genes) { 
    var baseURL = "http://"+window.location.hostname+":8080/datastore/genes/";
    url = baseURL + genes + "/avg"

    var avg
    $.ajax({
        async: false,
        cache: false,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            avg = JSON.parse(data);
        }
    });

    return avg

}



// Fetches std from datastore
function Std(gene) {
    
    var baseURL = "http://"+window.location.hostname+":8080/datastore/gene/";
    url = baseURL + gene + "/stddev"

    var res
    $.ajax({
        async: false,
        cache: false,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            res = data;
        }
    });

    //console.log(gene,"got std:", res)

    return res;

}


// Fetches std from datastore
function Var(gene) {
    
    var baseURL = "http://"+window.location.hostname+":8080/datastore/gene/";
    url = baseURL + gene + "/vari"

    var res
    $.ajax({
        async: false,
        cache: false,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            res = data;
        }
    });

    //console.log(gene,"got variance:", res)

    return res;

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

    //console.log("name: ", name)

    return name;

}



function GetBg(geneId,exprs) {
    var baseURL = "http://"+window.location.hostname+":8080/datastore/gene/"
        url = baseURL + geneId+"/"+exprs+"/bg"

    var info
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
             info = data;
        }
    });

    //console.log("Number of common genes: ", num)

    return info;

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

    //console.log("Number of common genes: ", num)

    return num;

}

function setScale(scale) { 
    console.log("js setting scale to ", scale);
    
    var baseURL = "http://"
                    +window.location.hostname
                    +":8080/datastore/setscale";
    
    updateColor(scale); 
        
    $.post( baseURL, scale, function( data ) {
            updateNodeColors();
            scaleDefer();
            
    });
    
    // erase content of infopanel
    document.getElementById('info-panel').innerHTML = '';



} 

function updateColor(scale) { 
    
    if(scale == "log") { 
        //console.log("log scale") 
        color = d3.scale.linear()
            .domain([-2,2])
            .range(colorbrewer.RdYlBu[3]);

    } 
    else { 
        //console.log("abs scale") 
        color = d3.scale.linear()
            .domain([-400,0,400])
            .range(colorbrewer.RdYlBu[3]);
    } 
} 




function visGenePanel(name){
    info = GetInfo(name);
    
    
    //console.log("The selected node was a gene!");

    // remove old info body
    document.getElementById('info-panel').innerHTML = '';

    // Set up new info box
    var panelDiv = document.createElement('div');
    panelDiv.className = 'panel panel-default';

    var panelHeadingDiv = document.createElement('div');
    panelHeadingDiv.id = 'info-panel-heading';
    panelHeadingDiv.className = 'panel-heading';
    var str = '<h5>'+info.Definition+'</h5>'
    panelHeadingDiv.innerHTML = str

    var panelBodyDiv = document.createElement('div');
    panelBodyDiv.id = 'info-panel-body';
    panelBodyDiv.className = 'panel-body';
    panelBodyDiv.innerHTML = GenerateInfoPanel(info)


    panelDiv.appendChild(panelHeadingDiv);
    panelDiv.appendChild(panelBodyDiv);

    document.getElementById('info-panel').appendChild(panelDiv);

    $(GetVis(info.Id)).appendTo(".visman"); 


} 
