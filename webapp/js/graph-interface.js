// The graph used for the visualization. added to own file
// to separate code for visualizing and keeping track of graph


var nodes = []; 
var edges = []; 
var graph; 


function Graph(cy){
    this.addNode = function(n){
        if(typeof findNode(n.id) != 'undefined') {
            return
        }

        n.graphics.name = n.graphics.name.split(" ")
        if(n.name === "\"bg\"") {
            //console.log(n.graphics.name[0])
            //var a = n.graphics.name[0].split("/")
            //var b = a[a.length-1]
            //var url = window.location.hostname+":8080/public/pathways/"+b 
            n.graphics.bgimage = n.graphics.name[0]
            n.graphics.bgcolor = "#fff"
            n.graphics.name = ""
        }
        else{
            n.graphics.bgimage = "";
        }

    
        n.graphics.name = n.graphics.name.toString().replace(/,/g, " ")
        
        n.graphics.valign = "center"
        if(n.graphics.name.substr(0,2) == "C0"){
            n.graphics.valign = "bottom"
        } 

        var no = {
            group: 'nodes',
            data: { 
                id: ''+ n.id,
                name: JSON.parse(n.name),
                graphics: n.graphics, 
                //background-image: "test.png",
            },
            position: {
                x: n.graphics.x,
                y: n.graphics.y
            },
            //grabbable: false,
            locked: true, 
        };
        var  gr = no.data.graphics; 

        if(gr.name == "" && gr.bgimage == ""){
            return
        } 
        nodes.push(no);
        cy.add(no); 
        
    };

    this.addEdge = function(e){
        var s = findNode(e.source); 
        var t = findNode(e.target); 
        
        if(typeof s == 'undefined' || typeof t == 'undefined'){
            return
        }

        var ed = {
            group: "edges",
            data: {
                source: ''+e.source,
                target: ''+e.target,
            },
        }; 

        edges.push(ed); 
        //cy.add(ed); 

    }

    var findNode = function (id) {
        for (var i in nodes) {
            if (nodes[i]["data"]["id"] === ''+id) {
                return nodes[i];
            }
        };
    };

    var update = function() {
        //cy.layout(); 
    }
}


