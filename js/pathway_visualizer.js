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
            name: 'random', 
            gravity: true,
            liveUpdate: true,
            maxSimulationtime: 1000,
        },
        
        showOverlay: false,
        minZoom: 0.5,
        maxZoom: 2,
        style: cytoscape.stylesheet()
            .selector('node')
            .css({
                'content': 'data(name)',
                'text-valign': 'right',
                'background-color': 'steelblue',
                'text-outline-width': 0,
                'text-outline-color': '#ccc',
                'text-opacity': 0.5,
                'text-color': '#ccc',
                'height': 10,
                'width': 10, 
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
                    prevSelection.cyTarget.css('background-color', 'steelblue');
                }
                d.cyTarget.css('background-color', '#2CA25F');
                prevSelection = d;

                // Determine selected node, can be gene/pathway/compound
                node = d.cyTarget.data();
                nodeType = node.name.split(":");


                info = GetInfo(d.cyTarget.data());


                if(nodeType[0] === 'hsa'){
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



                

        



                // write some contents to it
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
    var str = '<table class="table" style="word-wrap: break-word;table-layout:fixed">';

    str += '<thead><tr><th style="width: 10%"></th><th style="width: 90%"></th></tr></thead>'
    str += '<tbody>'
    str += '<tr><td>Expression:</td><td><div class="visman"></div></td></tr>';
    str += '<tr><td>Id:</td><td>hsa:' + info.Id + '</td><td>'
    str += '<tr><td>Definition:</td><td>' + info.Definition + '</td><td>'
    str += '<tr><td>Orthology:</td><td>' + info.Orthology + '</td><td>'
    str += '<tr><td>Organism:</td><td>' + info.Organism + '</td><td>'
    str += '<tr><td>Pathways:</td><td>' + info.Pathways + '</td><td>'
    str += '<tr><td>Diseases:</td><td>' + info.Diseases + '</td><td>'
    str += '<tr><td>Modules:</td><td>' + info.Modules + '</td><td>'
    str += '<tr><td>Drug target:</td><td>' + info.Drug_Target + '</td><td>'
    str += '<tr><td>Classes:</td><td>' + info.Classes + '</td><td>'
    str += '<tr><td>Position:</td><td>' + info.Position + '</td><td>'
    str += '<tr><td>Motif:</td><td>' + info.Motif + '</td><td>'
    str += '<tr><td>DB Links:</td><td>' + info.DBLinks + '</td><td>'
    str += '<tr><td>Structure:</td><td>' + info.Structure + '</td><td>'
    str += '<tr><td>AASeq:</td><td>' + info.AASEQ.Sequence + '</td><td>'
    str += '<tr><td>NTSeq:</td><td>' + info.NTSEQ.Sequence + '</td><td>'
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



