var wsURL
window.onload = function() {
    serverAddr = getVisServerAddress();
    wsURL = "ws://"+serverAddr; 
    console.log("visualization server is at:", wsURL); 
    console.log("Starting visualization..."); 
    loadCy(); 
};

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
            //name: 'arbor', 
           // gravity: true,
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

// Adding custom css to page 
function addCSS(cssPath) {
    linkElement = document.createElement("link");
    linkElement.rel = "stylesheet";
    linkElement.href = cssPath; 

    document.head.appendChild(linkElement);
}
addCSS("/css/pathway-visualizer.css"); 

