// The graph used for the visualization. added to own file
// to separate code for visualizing and keeping track of graph


var nodes = []; 
var edges = []; 
var graph; 

function Graph(cy){
    this.addNode = function(n){
        
        if(typeof findNode(n.id) != 'undefined') {
            console.log("attempted to add node (",n.id,") which exists.."); 
            return
        }

        var no = {
            group: 'nodes',
            data: { 
                id: ''+ n.id,
                name: JSON.parse(n.name),
                weight: 10,
                height: 10,
            },
            position: {
                x: Math.random() * 100,
                y: Math.random() * 100
            }
        };

        nodes.push(no);
        cy.add(no); 
    };

    this.addEdge = function(e){
        var s = findNode(e.source); 
        var t = findNode(e.target); 
        
        if(typeof s == 'undefined' || typeof t == 'undefined'){
            console.log("Attempted to add a faulty edge"); 
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
        cy.add(ed); 

    }

    var findNode = function (id) {
        for (var i in nodes) {
            if (nodes[i]["data"]["id"] === ''+id) {
                return nodes[i];
            }
        };
    };

    var update = function() {
        cy.layout(); 
    }
}


