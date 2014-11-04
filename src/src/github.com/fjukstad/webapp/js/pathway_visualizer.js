var wsURL
window.onload = function() {
    serverAddr = getVisServerAddress();
    wsURL = "ws://"+serverAddr; 
    loadCy(); 


    updateSettingsView()


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
       boxSelectionEnabled: false,
        style: cytoscape.stylesheet()
            .selector('node')
            .css({
                'content': 'data(graphics.name)',
                'text-valign': 'data(graphics.valign)',
                'background-color': 'data(graphics.bgcolor)',
                'background-image': 'data(graphics.bgimage)',
                'border-color': 'data(graphics.fgcolor)',
                'border-opacity': '1',
                'border-width': '2',
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
                    resizeHeader();
                }
                if(nodeType[0] === 'path'){
                    var pathid = node.name.split(":")[1]

                    var a = window.location.href
                    var b = a.split("=")
                    var c = b[0]
                    var url = c+"="+pathid

                    window.location.assign(url)

                }
                
                if(nodeType[0] === 'cpd'){  
                    console.log("Compound inbound", node) 
                    
                    visCompoundPanel(node.name); 
                    resizeHeader(); 

                }
                
                    
                d.cyTarget.edges().css({
                    'line-color': 'red'
                });


            });


            // Load data from JSON 
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
            hsas.push(name.split(" ")[0]);
        }
    }

    
    // convert list to string
    var hsastring = hsas.toString().replace(/,/g,"+")

    var foldchange = GetFoldChange(hsastring)
    var pvalues = GetPValues(hsastring) 

    // check if avg diff response is valid
    if (typeof foldchange === 'undefined'){
        alert("Unexpected error. Please try to refresh the web page")
    }; 

    console.log("FUCKERS", foldchange.Result) 

    var graphNodes = cy.nodes();

    for (var n in graphNodes) {
        if(n < graphNodes.length){
            if(graphNodes[n].style().shape == "rectangle"){
                var name = graphNodes[n].data().name.split(" ")[0];
                name = name.split(":")[1]
                console.log("NAME:", name) 
                fc = foldchange.Result[name];
                console.log("FC:",fc) 
                var c; 
                if(fc === "NA") { 
                   c = "#ffffff"
                } 
                else if(fc === undefined){
                    c = "#ffffff"
                } else { 
                    c = color(fc)
                }  
                graphNodes[n].css("background-color", c)
                
                p = pvalues.Result[name]
                if(p === "NA" || p === undefined){
                    c = "#000";
                    graphNodes[n].css("border-width", 1) 
                } else {
                    console.log("p-value:", p) 
                    if(parseFloat(p) < 0.009978){
                        c = "#67a9cf"

                    } else {
                        c = "#ef8a62" 
                    }
                }
                graphNodes[n].css("border-color", c) 

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
    var cyt = $('#cy')[0]
    cyt.style.height = $(window).height()-100+"px"
    cyt.style.width = $(".col-sm-8")[0].clientWidth-25+"px"
    
    // center the cytoscape graph after resize
    try { 
        cy.center() 
    } catch(TypeError){ 
    } 

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

        var panel = $("#info-panel-body")[0]
        //panel.style.height =  $(window).height()-100+"px"
        panel.style.width = $(".col-sm-4")[0].style.width;
        //clientWidth-50+"px"
    } catch(TypeError){ 
    } 
};


function resizeHeader(){
    try { 
        var header = $("#info-panel-heading")[0]
        header.style.width = $(".col-sm-4")[0].clientWidth-46+"px"

        var jmolview = $("img#jmolview")[0]
        jmolview.style.width = $("#c3").width()/2 + "px"


    } catch(TypeError){
    }
}

var settings =  {Smoking: true, HormoneTherapy: true, Disable: true} 

function updateSettingsView() {
    settings = getSettings() 
    console.log(settings) 
    // check checkboxes
    $( "input#disable").prop('checked', settings.Disable)
    $( "input#smoking").prop('checked', settings.Smoking)
    $( "input#hormones").prop('checked', settings.HormoneTherapy)
} 
