var wsURL
window.onload = function() {
    serverAddr = getVisServerAddress();
    wsURL = "ws://"+serverAddr; 
    console.log("visualization server is at:", wsURL); 
    console.log("Starting visualization..."); 
    loadCy(); 
};

var prevSelection;



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
                'text-valign': 'center',
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
            console.log("ready");
            graph = new Graph(cy); 

            cy.on('select', 'node', function(d){

                // Determine selected node, can be gene/pathway/compound
                node = d.cyTarget.data();
                nodeType = node.name.split(":");
                
                Pace.restart()
                

                if(nodeType[0] === 'hsa'){
                    
                    info = GetInfo(d.cyTarget.data());
                    
                    console.log("The selected node was a gene!");

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

                    $(GetVis(info.Id)).appendTo(".visman"); 


                
                }

                if(nodeType[0] === 'path'){
                    console.log("The selected node was a pathway!");
                }
                if(nodeType[0] === 'cpd'){
                    console.log("The selected node was a compund!");
                }
                    
                d.cyTarget.edges().css({
                    'line-color': 'red'
                });


            });


            /*
            cy.on('unselect', 'node', function(d){
               // d.cyTarget.css('background-color', 'steelblue');
            });

            cy.on('mouseup', '', function(d) {
            });
            */
    
            /*
            cy.on('zoom', function(d){
                var zoomLevel = cy.zoom();
                });
            */


            // Load data from JSON 
            console.log(wsURL)
            var socket = new WebSocket(wsURL); 
            socket.onmessage = function(m){
                var message = JSON.parse(m.data); 
                if(message.command == "\"InitGraph\""){
                    
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
                                                //graph.push(ed); 
                    }
                    cy.layout();
                }

                if(message.command == "\"AddNode\""){
                    graph.addNode(message); 
                    cy.layout();
                }
            } 
            
            
        }
    }; 


    $('#cy').cytoscape(options); 
/*
    $('#cy').cytoscapePanzoom({
        zoomFactor: 0.05, // zoom factor per zoom tick
        zoomDelay: 45, // how many ms between zoom ticks
        minZoom: 0.1, // min zoom level
        maxZoom: 10, // max zoom level
        fitPadding: 50, // padding when fitting
        panSpeed: 10, // how many ms in between pan ticks
        panDistance: 10, // max pan distance per tick
        panDragAreaSize: 75, // the length of the pan drag box in which the vector for panning is calculated (bigger = finer control of pan speed and direction)
        panMinPercentSpeed: 0.25, // the slowest speed we can pan by (as a percent of panSpeed)
        panInactiveArea: 8, // radius of inactive area in pan drag box
        panIndicatorMinOpacity: 0.5, // min opacity of pan indicator (the draggable nib); scales from this to 1.0
        autodisableForMobile: true, // disable the panzoom completely for mobile (since we don't really need it with gestures like pinch to zoom)

        // icon class names
        sliderHandleIcon: 'fa fa-minus',
        zoomInIcon: 'fa fa-plus',
        zoomOutIcon: 'fa fa-minus',
        resetIcon: 'fa fa-expand'
    });
    */
    
}


function GenerateInfoPanel(info){

    pathwayLinks = CreatePathwayLinks(info.Pathways)


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
    str += '<tr><td>Definition:</td><td>' + info.Definition + '</td><td>'
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

    // Fetch new expression values and colors for every gene
    var nodes = cy.nodes();
    for (var n in nodes) {
        if(n < nodes.length){
            if(nodes[n].style().shape == "rectangle"){
                var name = nodes[n].data().name.split(" ")[0];
                var c = color(AvgDiff(name))
                console.log(c)
                nodes[n].css("background-color", c)
            }
        }
    }
}

functionman = function() {
    console.log("updating color");
}



function savePathway()
{
    // get cytoscape instance
    var cy = $('#cy').cytoscape('get')
    // set image source
    $('#image')[0].src = cy.png()

} 


