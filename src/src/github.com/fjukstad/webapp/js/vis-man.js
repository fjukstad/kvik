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
        },
        error: function(xhr, ajaxOptions, thrownError){
            console.log("Bar chart not implemented yet");
            res = "";
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

function FoldChange(genes) { 
    var baseURL = "http://"+window.location.hostname+":8080/datastore/fc/";
    url = baseURL + genes

    var fc
    $.ajax({
        async: false,
        cache: false,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            fc = JSON.parse(data);
        }
    });

    console.log(fc) 

    return fc

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
        },
        error: function(xhr, ajaxOptions, thrownError){
            console.log("Standard deviation not implemented yet");
            res = "";
        }
    });


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
        },
        error: function(xhr, ajaxOptions, thrownError){
            console.log("Variance not implemented yet");
            res = "" 
        }
    });


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


    return name;

}



function GetBg(geneId,exprs) {
    /*
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
	*/
var info = "hepp"
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


    return num;

}

function getSettings(){
    var baseURL = "http://"
                    +window.location.hostname
                    +":8080/datastore/getsettings/all";
    var res
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: baseURL,
        dataType: "text",
        success: function(data) {
             res = data;
        },
        error: function(xhr, ajaxOptions, thrownError){
            console.log("Settings not implemented yet...");
            //res = "{Smoking: true, HormoneTherapy: false, Disable: true}"
            //return "" 
        }
    });

    return {Smoking:true, HormoneTherapy: false, Disable: true}
    //$.parseJSON(res) 


} 

function setSettings(smoking, hormones, disable) {
    var baseURL = "http://"
                    +window.location.hostname
                    +":8080/datastore/setsettings/";

    console.log("Setting settings:", settings) 
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

var colmax = 1.4//500//10000
var colmin = -0.7//-500//-1000

var colmaxlog = 1.4//20
var colminlog = -0.7

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
function visGenePanel(name){
    
    latestGene = name;
    if(latestGene === undefined){
        return
    }


    console.log(name)

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
    panelBodyDiv.innerHTML = GenerateInfoPanel(info)


    panelDiv.appendChild(panelHeadingDiv);
    panelDiv.appendChild(panelBodyDiv);

    document.getElementById('info-panel').appendChild(panelDiv);

    var viscode = GetVis(info.Id)
    console.log(viscode) 
    $(viscode).appendTo(".visman"); 


} 

function visCompoundPanel(name) {

    console.log(name)
    
    var info = GetInfo(name); 
    console.log(info) 

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
