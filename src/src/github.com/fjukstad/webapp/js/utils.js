
var colmax = 1.4//500//10000
var colmin = -0.7//-500//-1000

var colmaxlog = 1.4//20
var colminlog = -0.7

    
var color = d3.scale.linear()
    .domain([colmax,0,colmin])
    .range(colorbrewer.RdYlBu[3]);

var pcolor = d3.scale.linear()
    .domain([0.0001,0.00915,0.010])
    .range(colorbrewer.RdYlBu[3]);
    //.range(["yellow", "white", "blue"]);


function getPathwayId(){ 
    var currentLocation = window.location;
    var path = currentLocation.pathname
    var pathwayid = path.split("=")[1]
    return pathwayid
} 
function updatePathwayInfoPanel(){

    var id = getPathwayId() 
    var p = GetPathway(id) 



    console.log("GOT THAT PATHWAY: ", p) 

    var div = d3.select("body").selectAll("#pathway-info-view"); 

    var pg = div.append("div")
        .attr("class", "panel-group")
        .attr("id", "pw-accord")

    var p1 = pg.append("div")
                .attr("class", "panel panel-default panel-collapse")

    var p1h = p1.append("div")
                .attr("class", "panel-heading")
                .append("h4")
                .attr("class", "panel-title") 
                .text(p.Name) 
                
    var descr = p.Description.split(".")[0] + " (read more) "

    if(descr == " (read more) ") {
        descr = "More information here"
    } 

    var p1b = p1.append("div") 
                .attr("class", "panel-body")
                .append("a")
                .attr("href", "http://www.kegg.jp/dbget-bin/www_bget?"+id)
                .attr("target","_blank")
                .text(descr) 
} 

function GenerateInfoPanel(info, genename){
    pathwayLinks = CreatePathwayLinks(info.Pathways)

    var std = parseFloat(Std(info.Id)).toFixed(3) 
    var variance = parseFloat(Var(info.Id)).toFixed(3)

    var keggid = "" 
    var lookupId = "" 
    if(genename != undefined){
        keggid = genename
        lookupId = genename
    } else { 
        keggid = "hsa:"+info.Id
        lookupId = info.Id 
    }
        
    var f = GetFoldChange(keggid).Result[lookupId]
    var fc = parseFloat(f).toFixed(3)

    var p = GetPValues(keggid).Result[lookupId]
    var pvalue = parseFloat(p).toFixed(5) 


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
    str += '<small>P-Value: '+pvalue+'</br>Fold Change: '+fc
    //str +='</br>Standard deviation: '+std+'</br>Variance:'+variance+ '</small>'
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
