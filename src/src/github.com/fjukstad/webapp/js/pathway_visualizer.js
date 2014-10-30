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
    str += 'Gene expression values for entire dataset' 
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
    str += 'Similar pathways'
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
    str += '<tr><td>Id:</td><td><a href="http://www.genome.jp/dbget-bin/www_bget?hsa:'+info.Id+'" target="_blank">hsa:' + info.Id + '</a></td><td>'
    str += '<tr><td>Definition:</td><td>' + info.Name + '</td><td>'
    str += '<tr><td>Orthology:</td><td>' + info.Orthology + '</td><td>'
    //str += '<tr><td>Organism:</td><td>' + info.Organism + '</td><td>'
    if(info.Diseases){
        str += '<tr><td>Diseases:</td><td>' + info.Diseases + '</td><td>'
    }
    if(info.Modules){ 
        str += '<tr><td>Modules:</td><td>' + info.Modules + '</td><td>'
    }
    if(info["Drug_Target"]){
        str += '<tr><td>Drug target:</td><td>' + info["Drug_Target"] + '</td><td>'
    } 
    str += '<tr><td>Classes:</td><td>' + info.Classes + '</td><td>'
    str += '<tr><td>Position:</td><td>' + info.Position + '</td><td>'
    str += '<tr><td>Motif:</td><td>' + info.Motif + '</td><td>'
    str += '<tr><td>DB Links:</td><td>' + CreateDBLinks(info.DBLinks) + '</td><td>'
    str += '<tr><td>Structure:</td><td>' + FetchJMOL(info.Structure) + '</td><td>'
    //str += '<tr><td>AASeq:</td><td>' + info.AASEQ.Sequence + '</td><td>'
    //str += '<tr><td>NTSeq:</td><td>' + info.NTSEQ.Sequence + '</td><td>'
    str += '</tbody>'
    str += '</table>';
    str += '</div></div></div>'

    console.log(info.DBLinks)

    str += '</div>'
    
       return str
}

function GenerateCompoundInfoPanel(info) {

    //http://www.genome.jp/Fig/compound/C00575.gif

    var str = '<div class="panel-group" id="accordion">'
    str += '<div class="panel panel-default">';
    str += '<div class="panel-heading">'
    str += '<h4 class="panel-title">'
    str += '<a data-toggle="collapse" data-parent="#accordion" href="#c1">'
    str += 'Structure'
    str += '</a> </div>'
    str += '<div id="c1" class="panel-collapse collapse in">'
    str += '<div class="panel-body">'
    // Fetch structure vis from kegg 
    var structURL = "http://www.genome.jp/Fig/compound/"+info.Entry+".gif" 
    str += '<div class="visman"><img src="'+structURL+'" class="structure"></img></div>'
    str += '</div></div></div>'

    pathwayLinks = CreatePathwayLinks(info.Pathway)
    str += '<div class="panel panel-default">';
    str += '<div class="panel-heading">'
    str += '<h4 class="panel-title">'
    str += '<a data-toggle="collapse" data-parent="#accordion" href="#c2">'
    str += 'Similar Pathways'
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
    str += '<div id="c3" class="panel-collapse">'
    str += '<div class="panel-body">'

    str += '<table class="table" style="word-wrap: break-word;table-layout:fixed">';
    str += '<thead><tr><th style="width: 20%"></th><th style="width: 80%"></th>'
    str += '<tbody>'
    str += '<tr><td>Id:</td><td><a href="http://www.genome.jp/dbget-bin/www_bget?'+info.Entry+'" target="_blank">' + info.Entry + '</a></td><td>'
    str += '<tr><td>Name:</td><td>' + CreateNameList(info.Name) + '</td><td>'
    str += '<tr><td>DB Links:</td><td>' + CreateCompoundDBLinks(info.DBLinks) + '</td><td>'



    return str
} 

function CreateNameList(names) { 
    var res = ""; 

    for(i in names){
        res += names[i]+"</br>"
    }
    
    return res
} 

function FetchJMOL(structure) {
    try { 
        var ids = structure.split(" ")
        var id = ids[1].toLowerCase()
        console.log(ids) 
        var link = "http://www.genome.jp/Fig/pdb/pdb"+id+".png"
        var res = '<a href="'+link+'" target="_blank"><img src="'+link+'" id="jmolview"></a>'
        return res
    } catch(TypeError) {
        return ""
    }
}

function CreateCompoundDBLinks(links) { 

    var res = ""; 
    if(links["3DMET"]){
        var tredmet = '<a href="http://www.3dmet.dna.affrc.go.jp/cgi/show_data.php?acc='+links["3DMET"]+'" target="_blank">3DMET</a>';
        res += tredmet + "</br>";
    }
    if(links["PubChem"]){
        var pubchem = '<a href="http://pubchem.ncbi.nlm.nih.gov/summary/summary.cgi?sid='+links["PubChem"]+'" target="_blank">PubChem</a>';
        res += pubchem + "</br>";
    }

    if(links["ChEBI"]){
        var ChEBI = '<a href="http://www.ebi.ac.uk/chebi/searchId.do?chebiId=CHEBI:'+links["ChEBI"]+'" target="_blank">ChEBI</a>'
        res += ChEBI + "</br>";
    } 

    
        
    return res 


} 

function CreateDBLinks(links) {
    
    var res = "" 
    try { 
    var gname = '<a href="http://www.genenames.org/cgi-bin/search?search_type=symbols&search='+links.HGNC+'" target="_blank">GeneNames</a>'
    res += gname + "</br>"
    
    var ensembl = '<a href="http://www.ensembl.org/Multi/Search/Results?q='+links.Ensembl+'" target="_blank">Ensembl</a>'

    res += ensembl + "</br>"

    var ncbigeneid = '<a href="http://www.ncbi.nlm.nih.gov/gene/?term='+links["NCBI-GeneID"]+'" target="_blank">NCBI Gene </a>'

    res += ncbigeneid + "</br>"


    var uniprot = '<a href="http://www.uniprot.org/uniprot/'+links.UniProt+'" target="_blank">UniProt</a>'

    console.log(uniprot) 
    res += uniprot

    } catch (TypeError){
        console.log(links);
        console.log(TypeError)
    }
    return res
    
} 


function CreatePathwayLinks(ids) {
    var baseURL = "http://"+window.location.hostname+":8000/browser/pathwaySelect="
    links  = "" 

    var currentLocation = window.location;
    var path = currentLocation.pathname
    var pathwayid = path.split("=")[1]
    for (i in ids) {
        id = ids[i];

        // We only care about human pathways. 
        id = id.replace("map", "hsa")
        
        if (id != pathwayid) {
            name = GetPathwayName(id)
            if (name != "") { 
                pathwayIds = id+"+"+pathwayid
                num = GetCommonGenes(pathwayIds)
                test = "<div style=\" float: right; display: inline-block; width:" 
                test += num
                test += "px; height: 10px; background-color: #a6bbc8\"></div>"

                links += "<a href=\""+baseURL+id+"\" title=\""+id+"\">"+name+"</a>"
                links += test + "</br>"
            }
            
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

    var foldchange = FoldChange(hsastring)

    // check if avg diff response is valid
    if (typeof foldchange === 'undefined'){
        alert("Unexpected error. Please try to refresh the web page")
    }; 

    var graphNodes = cy.nodes();

    for (var n in graphNodes) {
        if(n < graphNodes.length){
            if(graphNodes[n].style().shape == "rectangle"){
                var name = graphNodes[n].data().name.split(" ")[0];
                fc = foldchange.Result[name]
                console.log("foldman", fc) 
                if(fc === "NA") { 
                   var c = "#ffffff"
                } else { 
                    var c = color(fc)
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
