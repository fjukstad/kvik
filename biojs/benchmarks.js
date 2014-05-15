

function PostResults(method,data) {
    var url = "http://localhost:4444/time/"+method
    $.post(url, data) 
}

function PostSamples(method,data){
    var url = "http://localhost:4444/samples/"+method
    $.post(url, data) 
} 

function LoadPathway(pathwayId) { 
    serverAddr = NewGraphServer(name);
    wsURL = "ws://"+serverAddr;

   $("#cy").innerHTML = ""; 
    
    loadCy(); 
    

} 

function foobar_cont(){
    console.log("finished.");
};

function sleep(millis, callback) {
    setTimeout(function()
            { callback(); }
    , millis);
}


function keggview() {
    var instance = new Biojs.KEGGViewer({
         target: 'YourOwnDivId',
         pathId: 'hsa04910',
         proxyUrl: 'proxy.php',
         expression:{
             upColor:'red',
             downColor:'blue',
             genes: ['hsa:7248', 'hsa:51763', 'hsa:2002', 'hsa:2194'],
             conditions: [
                 {
                     name: 'condition 1',
                     values: [-1, 0.5, 0.7, -0.3]
                 },
                 {
                     name: 'condition 2',
                     values: [0.5, -0.1, -0.2, 1]
                 },
                 {
                     name: 'condition 3',
                     values: [0, 0.4, -0.2, 0.5]
                 }
               ]
         }     
     
    });
     
};

// current pathway
var logpwid = window.location.pathname.split("=")[1]; 

var numRuns = 0
function StartBenchmarks(){ 

    minSamples = 200

    var suite = new Benchmark.Suite('test', { 
        'onComplete': function(){
            for(i = 0; i < suite.length; i++){
                bench = suite[i];
                console.log(bench)
                
                // first add some more numbers to the benchmark results
                bench.stats.count = bench.count;
                bench.stats.hz = bench.hz;
                bench.stats.timeElapsed = bench.times.elapsed;
                
                // send results to server
                PostResults(bench.name, JSON.stringify(bench.stats))
            }

        }
    });

    

        
    suite.add("keggview", function(deferred){
            keggview()
            defff = deferred
            numRuns = numRuns + 1 
            console.log(numRuns)
        },
        { 
            'defer': true,
            'minSamples': minSamples,
            'onComplete': function(e){
                console.log(e.currentTarget.name," is done running")
                numRuns = 0
            }
        } 
    );

    suite.run({'async':true});
    benchmarked = true

} 



var defff

function deferAway(){ 
    
    try {
        defff.resolve()
        nodes = []; 
        edges = []; 
        $("#YourOwnDivId").empty()
    }
    catch(e){
        console.log(e)
    } 
} 

function scaleDefer(){
    try {
        defff.resolve();
    }
    catch(e){
        return 
    }
} 

