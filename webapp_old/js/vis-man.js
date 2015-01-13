
function GET(url, parse) { 
    var response
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            if(parse) {
                response = JSON.parse(data);
            } else {
                response = data
            } 
        },
        error: function(xhr, ajaxOptions, thrownError){
            response = "";
        }
    });
    return response
}

function GetPathwayGraph(id) {
    var url = "http://"+window.location.hostname+":8080/pathwayGraph/"+id
    return GET(url, true) 
}

// js for retrieving vis about gene
function GetVis(gene) {
    var baseURL = "http://"+window.location.hostname+":8080/vis/";
    url = baseURL + gene
    return GET(url,false) 
}

// Will download a parallel coordinate plot 
function GetParallelVis() {
    var baseURL = "http://"+window.location.hostname+":8080/parallel";
    url = baseURL
    return GET(url, false) 
}

function GetFoldChange(genes) { 
    var baseURL = "http://"+window.location.hostname+":8080/datastore/fc/";
    url = baseURL + genes
    return GET(url,true) 
}

// Fetches std from datastore
function Std(gene) {
    var baseURL = "http://"+window.location.hostname+":8080/datastore/gene/";
    url = baseURL + gene + "/stddev"
    return GET(url,false) 
}


// Fetches std from datastore
function Var(gene) {
    var baseURL = "http://"+window.location.hostname+":8080/datastore/gene/";
    url = baseURL + gene + "/vari"
    return GET(url, false) 
}



function GetPathwayName(id) {
    var baseURL = "http://"+window.location.hostname+":8080/info/pathway/"
    url =  baseURL+id+"/name"
    return GET(url, false) 
}

function GetPathway(id) { 
    var baseURL = "http://"+window.location.hostname+":8080/info/"
    url =  baseURL+id+"/all"
    return GET(url, true) 

}

function GetBg(geneId,exprs) {
    var baseURL = "http://"+window.location.hostname+":8080/datastore/gene/"
        url = baseURL + geneId+"/"+exprs+"/bg"
    return GET(url, false) 
}

function GetCommonGenes(ids) {
    var baseURL = "http://"+window.location.hostname+":8080/info/pathway/"
        url =  baseURL+ids+"/commongenes"
    return GET(url, false) 
}

function getSettings(){
    var url = "http://"
                    +window.location.hostname
                    +":8080/datastore/getsettings/all";
    //return GET(url, true) 

    return {Smoking:true, HormoneTherapy: false, Disable: true}


} 

function GetPValues(genes){
    var baseURL = "http://"+window.location.hostname+":8080/datastore/pvalues/";
    url = baseURL + genes
    return GET(url, true) 
}

function GetGenes(){
    var url = "http://"+window.location.hostname+":8080/datastore/genes/all";
    return GET(url, true) 
} 

function GetGeneId(genename) {
    var url = "http://"+window.location.hostname+":8080/geneid/"+genename
    return GET(url, false) 
}


// js for retrieving info about specific pathway/gene/compund etc.
function GetInfo(name) {
    var a = name;
    var b = a.replace(/ /g, '+').toLowerCase();

    var baseURL = "http://"+window.location.hostname+":8080/info";
    var infoType = "all";

    url = baseURL + "/" + b + "/" + infoType;
    return GET(url, true) 

}


function setSettings(smoking, hormones, disable) {
    var baseURL = "http://"
                    +window.location.hostname
                    +":8080/datastore/setsettings/";

    var s = settings// {Smoking: smoking, HormoneTherapy: hormones, Disable: disable} 
    
    $.post( baseURL, JSON.stringify(s));
    updateNodeColors();
    scaleDefer();
    visGenePanel(latestGene)

} 

function setScale(scale) { 

    var baseURL = "http://"
                    +window.location.hostname
                    +":8080/datastore/setscale";
    
    updateColor(scale); 
        
    $.post( baseURL, scale, function( data ) {
            updateNodeColors();
            scaleDefer();
            
    });
    
    visGenePanel(latestGene)


} 



function updateColor(scale) { 
    
    if(scale == "log") { 
        color = d3.scale.linear()
            .domain([colmaxlog,0,colminlog])
            .range(colorbrewer.RdYlBu[3]);

    } 
    else { 
        color = d3.scale.linear()
            .domain([colmax,0, colmin])
            .range(colorbrewer.RdYlBu[3]);
    } 
} 


// Latest gene variable used for resizing of window
var latestGene;  
function visGenePanel(name,genename){
    
    latestGene = name;
    if(latestGene === undefined){
        return
    }



    info = GetInfo(name);
    
    
    // remove old info body
    document.getElementById('info-panel').innerHTML = '';

    // Set up new info box
    var panelDiv = document.createElement('div');
    panelDiv.className = 'panel panel-default';

    var panelHeadingDiv = document.createElement('div');
    panelHeadingDiv.id = 'info-panel-heading';
    panelHeadingDiv.className = 'panel-heading';
    var str = '<h5>'+info.Name+'</h5>'
    panelHeadingDiv.innerHTML = str

    var panelBodyDiv = document.createElement('div');
    panelBodyDiv.id = 'info-panel-body';
    panelBodyDiv.className = 'panel-body';
    panelBodyDiv.innerHTML = GenerateInfoPanel(info,genename)

    panelDiv.appendChild(panelHeadingDiv);
    panelDiv.appendChild(panelBodyDiv);

    document.getElementById('info-panel').appendChild(panelDiv);
    
    if(genename == undefined) {
        genename = info.Name.split(",")[0]
    }
        

    var viscode = GetVis(genename)
    $(viscode).appendTo(".visman"); 

} 

function visCompoundPanel(name) {

    
    var info = GetInfo(name); 

    // remove old info body
    document.getElementById('info-panel').innerHTML = '';

    // Set up new info box
    var panelDiv = document.createElement('div');
    panelDiv.className = 'panel panel-default';

    var panelHeadingDiv = document.createElement('div');
    panelHeadingDiv.id = 'info-panel-heading';
    panelHeadingDiv.className = 'panel-heading';
    var str = '<h5>'+info.Name[0]+'</h5>'
    panelHeadingDiv.innerHTML = str

    var panelBodyDiv = document.createElement('div');
    panelBodyDiv.id = 'info-panel-body';
    panelBodyDiv.className = 'panel-body';
    panelBodyDiv.innerHTML = GenerateCompoundInfoPanel(info)
    
    panelDiv.appendChild(panelHeadingDiv);
    panelDiv.appendChild(panelBodyDiv);

    document.getElementById('info-panel').appendChild(panelDiv);

} 
