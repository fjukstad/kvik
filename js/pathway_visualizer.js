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



$(loadCy = function(){
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
                'background-color': '#000',
                'line-color': '#000',
                'target-arrow-color': '#000',
                'text-outline-color': '#000'
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
                    

            cy.on('mouseover', 'node', function(d){
                // update visuals of nodes
                if(prevSelection !== undefined){
                    //prevSelection.cyTarget.css('background-color', 'steelblue');
                    //prevSelection.cyTarget.css('text-opacity', '0.0');
                }
                //d.cyTarget.css('background-color', '#2CA25F');
                //d.cyTarget.css('text-opacity', '0.5');
                d.cyTarget.css({
                    'height': 'data(graphics.height + 50)',
                    'width': 'data(graphics.width + 50) ', 
                })



                prevSelection = d;
            });

            cy.on('select', 'node', function(d){

                // Determine selected node, can be gene/pathway/compound
                node = d.cyTarget.data();
                nodeType = node.name.split(":");


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
                    
                
                console.log("Neighbors: ", d.cyTarget.edges());

                d.cyTarget.edges().css({
                    'line-color': 'red'
                });


            });



            cy.on('select', 'node', function(d){
                //d.cyTarget.css('background-color', '#FEC44F');

                
                if(d.type === 'select') { 
                    /*
                    console.log("highlighted nodes: ", cy.elements("node:selected"));

                    // remove old info body
                    document.getElementById('info-panel').innerHTML = '';
        
                    // Set up new info box
                    var panelDiv = document.createElement('div');
                    panelDiv.className = 'panel panel-default';
        
                    var panelHeadingDiv = document.createElement('div');
                    panelHeadingDiv.id = 'info-panel-heading';
                    panelHeadingDiv.className = 'panel-heading';
                    var str = '<h5> Time series for selected genes </h5>'
                    panelHeadingDiv.innerHTML = str

                    var panelBodyDiv = document.createElement('div');
                    panelBodyDiv.id = 'info-panel-body';
                    panelBodyDiv.className = 'panel-body';
                    panelBodyDiv.innerHTML = GenerateParallelPanel()


                    panelDiv.appendChild(panelHeadingDiv);
                    panelDiv.appendChild(panelBodyDiv);

                    document.getElementById('info-panel').appendChild(panelDiv);

                    $(GetParallelVis()).appendTo(".parallel"); 
                    */
                }

                
            });

            cy.on('unselect', 'node', function(d){
               // d.cyTarget.css('background-color', 'steelblue');
            });

            cy.on('mouseup', '', function(d) {
            });
    

            cy.on('zoom', function(d){
                var zoomLevel = cy.zoom();
                /*
                if(zoomLevel >= 1.5){
                    cy.nodes().animate({
                      css: { 'text-opacity': '0.5' } }
                    , {
                      duration: 0
                    });
                }
                else {
                    cy.nodes().animate({
                      css: { 'text-opacity': '0.5' } }
                    , {
                      duration: 0
                    });
                }

                */
                });


            // Load data from JSON 
            var socket = new WebSocket(wsURL); 
            socket.onmessage = function(m){
                var message = JSON.parse(m.data); 
                if(message.command == "\"InitGraph\""){
                    
                    json = JSON.parse(JSON.parse(message.graph)); 
                    var numAdded = 0; 
                    console.log(json.nodes);
                    
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

});


function GenerateInfoPanel(info){

    pathwayLinks = CreatePathwayLinks(info.Pathways)


    var str = '<table class="table" style="word-wrap: break-word;table-layout:fixed">';
    str += '<thead><tr><th style="width: 20%"></th><th style="width: 80%"></th></tr></thead>'
    str += '<tbody>'
    str += '<tr><td>Expression:</td><td><div class="visman"></div></td></tr>';
    str += '<tr><td>Id:</td><td>hsa:' + info.Id + '</td><td>'
    str += '<tr><td>Definition:</td><td>' + info.Definition + '</td><td>'
    str += '<tr><td>Orthology:</td><td>' + info.Orthology + '</td><td>'
    //str += '<tr><td>Organism:</td><td>' + info.Organism + '</td><td>'
    str += '<tr><td>Pathways:</td><td>' + pathwayLinks + '</td><td>'
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
    return str
}

function CreatePathwayLinks(ids) {
    var baseURL = "http://"+window.location.hostname+":8000/demo/pathwaySelect="
    links  = "" 
    console.log(ids)
    for (i in ids) {
        id = ids[i];
        console.log(ids[i]);
        name = GetPathwayName(id)
        console.log(name)
        links += "<a href=\""+baseURL+id+"\" title=\""+id+"\">"+name+"</a></br>"
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



