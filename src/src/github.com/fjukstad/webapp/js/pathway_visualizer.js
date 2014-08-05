var wsURL
window.onload = function() {
    serverAddr = getVisServerAddress();
    wsURL = "ws://"+serverAddr; 
    //console.log("visualization server is at:", wsURL); 
    //console.log("Starting visualization..."); 
    loadCy(); 
};

var prevSelection;

var benchmarked = false
var drawnPathway = false

function getVisServerAddress() {
    var baseURL = "http://"+window.location.hostname+":8080"
    var visType = "/new/graph/pathway/"
    var pathwayIndex = 4; 
    var selectedPathways = window.location.href.split('/')[pathwayIndex];
    var url = baseURL+visType+selectedPathways;
    var serverURL; 

    $.ajax({
        async: false,
        cache: false,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data){
            serverURL = window.location.hostname+data; 
        }
    }); 
    return serverURL;
}



loadCy = function(){
    options = {
        layout: {
            name: 'preset', 
            fit: true,
        },
        
        showOverlay: false,
        minZoom: 0.2,
        maxZoom: 5,
        style: cytoscape.stylesheet()
            .selector('node')
            .css({
                'content': 'data(graphics.name)',
                'text-valign': 'data(graphics.valign)',
                'background-color': 'data(graphics.bgcolor)',
                'background-image': 'data(graphics.bgimage)',
                'border-color': 'data(graphics.fgcolor)',
                'border-opacity': '1',
                'border-width': '1',
                'text-outline-width': '0',
                'text-outline-color': '#fff',
                'text-opacity': 0.9,
                'color': '#000',
                'shape': 'data(graphics.shape)',
                'height': 'data(graphics.height)',
                'width': 'data(graphics.width)', 
                'font-family': 'helvetica',
                'font-size': 10,
            })
            .selector(':selected')
            .css({
                /*
                'background-color': '#000',
                'line-color': '#000',
                'target-arrow-color': '#000',
                'text-outline-color': '#000'
                */
            })
            .selector('edge')
            .css({
                'target-arrow-shape': 'triangle'
        }),
        elements : {
            nodes: [],
            edges : []
        },

        ready: function(){
            cy = this;
            graph = new Graph(cy); 

            drawnPathway = false

            cy.on('select', 'node', function(d){

                // Determine selected node, can be gene/pathway/compound
                node = d.cyTarget.data();
                nodeType = node.name.split(":");
                
                Pace.restart()

                if(nodeType[0] === 'hsa'){
                    var name = d.cyTarget.data().name
                    visGenePanel(name)
                }
                /*
                if(nodeType[0] === 'path'){
                    continue
                }
                if(nodeType[0] === 'cpd'){
                    continue
                }
                */
                    
                d.cyTarget.edges().css({
                    'line-color': 'red'
                });


            });


            // Load data from JSON 
            //console.log(wsURL)
            var socket = new WebSocket(wsURL); 
            socket.onmessage = function(m){
                var message = JSON.parse(m.data); 
                if(message.command == "\"InitGraph\""){
                    
                    //console.log("INIT ME")
                    json = JSON.parse(JSON.parse(message.graph)); 
                    var numAdded = 0; 
                    
                    for(var i in json.nodes){
                        var n = json.nodes[i]; 
                        graph.addNode(n); 
                    }
                    //cy.layout(); 
                    var cy_nodes = cy.add(nodes); 
                    for(var j in json.edges){
                        var e = json.edges[j]; 
                        graph.addEdge(e); 
                    }
                    cy.layout();

                    drawnPathway = true   

                    if(!benchmarked){ 
                        StartBenchmarks()
                        benchmarked = true
                    } 

                    updateNodeColors()

                    // WARNING: CLOSING SOCKET AFTER INIT
                    socket.close()

                    deferAway() 

                }

                if(message.command == "\"AddNode\""){
                    graph.addNode(message); 
                    cy.layout();
                }
            } 
            
            
        },
    }; 
    
    $('#cy').cytoscape(options); 
    resizeViews() 
}


function GenerateInfoPanel(info){
    pathwayLinks = CreatePathwayLinks(info.Pathways)


    var std = parseFloat(Std(info.Id)).toFixed(3) 
    var variance = parseFloat(Var(info.Id)).toFixed(3)
    var mean = parseFloat(AvgDiff(info.Id)).toFixed(3)


    var str = '<div class="panel-group" id="accordion">'
        
    str += '<div class="panel panel-default">';
    str += '<div class="panel-heading">'
    str += '<h4 class="panel-title">'
    str += '<a data-toggle="collapse" data-parent="#accordion" href="#c1">'
    str += 'Expression'
    str += '</a> </div>'
    str += '<div id="c1" class="panel-collapse collapse in">'
    str += '<div class="panel-body">'
    str += '<div class="visman"></div>'
    //str += '<button id="sort" onclick="sortBars()">Sort</button>'
    str += '<small>Mean: '+mean+'</br>Standard deviation: '+std+'</br>Variance:'+variance+ '</small>'
    str += '<div id="dsidinfo"></div>'
    str += '</div></div></div>'


    str += '<div class="panel panel-default">';
    str += '<div class="panel-heading">'
    str += '<h4 class="panel-title">'
    str += '<a data-toggle="collapse" data-parent="#accordion" href="#c2">'
    str += 'Pathways'
    str += '</a> </div>'
    str += '<div id="c2" class="panel-collapse collapse in">'
    str += '<div class="panel-body">'
    str += pathwayLinks
    str += '</div></div></div>'

    str += '<div class="panel panel-default">';
    str += '<div class="panel-heading">'
    str += '<h4 class="panel-title">'
    str += '<a data-toggle="collapse" data-parent="#accordion" href="#c3">'
    str += 'More information'
    str += '</a> </div>'
    str += '<div id="c3" class="panel-collapse collapse">'
    str += '<div class="panel-body">'

    str += '<table class="table" style="word-wrap: break-word;table-layout:fixed">';
    str += '<thead><tr><th style="width: 20%"></th><th style="width: 80%"></th>'
    str += '<tbody>'
    str += '<tr><td>Id:</td><td>hsa:' + info.Id + '</td><td>'
    str += '<tr><td>Definition:</td><td>' + info.Name + '</td><td>'
    str += '<tr><td>Orthology:</td><td>' + info.Orthology + '</td><td>'
    //str += '<tr><td>Organism:</td><td>' + info.Organism + '</td><td>'
    str += '<tr><td>Diseases:</td><td>' + info.Diseases + '</td><td>'
    str += '<tr><td>Modules:</td><td>' + info.Modules + '</td><td>'
    str += '<tr><td>Drug target:</td><td>' + info.Drug_Target + '</td><td>'
    str += '<tr><td>Classes:</td><td>' + info.Classes + '</td><td>'
    str += '<tr><td>Position:</td><td>' + info.Position + '</td><td>'
    str += '<tr><td>Motif:</td><td>' + info.Motif + '</td><td>'
    str += '<tr><td>DB Links:</td><td>' + info.DBLinks + '</td><td>'
    str += '<tr><td>Structure:</td><td>' + info.Structure + '</td><td>'
    //str += '<tr><td>AASeq:</td><td>' + info.AASEQ.Sequence + '</td><td>'
    //str += '<tr><td>NTSeq:</td><td>' + info.NTSEQ.Sequence + '</td><td>'
    str += '</tbody>'
    str += '</table>';
    str += '</div></div></div>'

    

    str += '</div>'
    
       return str
}

function CreatePathwayLinks(ids) {
    var baseURL = "http://"+window.location.hostname+":8000/demo/pathwaySelect="
    links  = "" 

    var currentLocation = window.location;
    var path = currentLocation.pathname
    var pathwayid = path.split("=")[1]
    for (i in ids) {
        id = ids[i];
        if (id != pathwayid) {
            name = GetPathwayName(id)
            pathwayIds = id+"+"+pathwayid
            num = GetCommonGenes(pathwayIds)
            test = "<div style=\" float: right; display: inline-block; width:" 
            test += num
            test += "px; height: 10px; background-color: #a6bbc8\"></div>"

            links += "<a href=\""+baseURL+id+"\" title=\""+id+"\">"+name+"</a>"
            links += test + "</br>"
            
        }
    }
    return links
} 

function GenerateParallelPanel() {
    var str = '<table class="table" style="word-wrap: break-word;table-layout:fixed">';
    str += '<thead><tr><th style="width: 20%"></th><th style="width: 80%"></th></tr></thead>'
    str += '<tbody>' 
    str += '<tr><td>Expression :</td><td><div class="parallel"></div></td></tr>';
    str += '</tbody>'
    str += '</table>';
    return str

}

// Adding custom css to page 
function addCSS(cssPath) {
    linkElement = document.createElement("link");
    linkElement.rel = "stylesheet";
    linkElement.href = cssPath; 

    document.head.appendChild(linkElement);
}
addCSS("/css/pathway-visualizer.css"); 


window.onerror = function(error) {
    alert(error);
};

function updateNodeColors() {


    // get list of genes in pathwaymap
    var hsas = [];
        
    for(i=0;i<nodes.length;i++){
        var n = nodes[i]; 
        name=n.data.name;
        if(!name.indexOf("hsa")){
            hsas.push(name.split(" "))[0];
        }
    }
    
    // convert list to string
    var hsastring = hsas.toString().replace(/,/g,"+")

    var ex = AvgDiffs(hsastring)

    // check if avg diff response is valid
    if (typeof ex === 'undefined'){
        alert("Unexpected error. Please try to refresh the web page")
    }; 

    var graphNodes = cy.nodes();

    for (var n in graphNodes) {
        if(n < graphNodes.length){
            if(graphNodes[n].style().shape == "rectangle"){
                var name = graphNodes[n].data().name.split(" ")[0];
                avg = ex.Expression[name]
                if(avg === 0) { 
                   var c = "#ffffff"
                } else { 
                    var c = color(avg)
                }  
                graphNodes[n].css("background-color", c)
            }
        }
    }
}



function savePathway()
{
    // get cytoscape instance
    var cy = $('#cy').cytoscape('get')
    // set image source
    $('#image')[0].src = cy.png()

} 

function ShowBgInfo(id,exprs) {

    var bg = GetBg(id,exprs);
    //bg = JSON.parse(GetBg(id,exprs));

    document.getElementById('dsidinfo').innerHTML =  bg

} 

window.onresize = function(event) {
    resizeViews();
} 

function resizeViews(){
    console.log("window is resizing") 
    var cy = $('#cy')[0]
    cy.style.height = $(window).height()-100+"px"
    cy.style.width = $(".col-sm-8")[0].clientWidth-25+"px"
    $('cy').resize() 


    console.log($(".col-sm-4"))
    $(".col-sm-4").client

      
    try { 
        // update gene panel
        visGenePanel(latestGene)
    } catch (TypeError) {
    }

    var exprs = $("#expression-view")[0]
    exprs.style.height =  $(window).height()-100+"px"
    exprs.style.width = $(".col-sm-4")[0].clientWidth-30+"px"
    
    try { 
        var header = $("#info-panel-heading")[0]
        //header.style.height =  $(window).height()-100+"px"
        header.style.width = $(".col-sm-4")[0].clientWidth-46+"px"
        console.log(header) 

        var panel = $("#info-panel-body")[0]
        //panel.style.height =  $(window).height()-100+"px"
        panel.style.width = $(".col-sm-4")[0].style.width;
        console.log("DICKBALLS", panel.style.width) 
        //clientWidth-50+"px"
    } catch(TypeError){ 
    } 





};



