

function PostResults(method,data) {
    var url = "http://localhost:4444/time/"+method
    $.post(url, data) 
}

function PostSamples(method,data){
    var url = "http://localhost:4444/samples/"+method
    $.post(url, data) 
} 

function NewGraphServer(name) { 
    var baseURL = "http://"+window.location.hostname+":8080"
    var visType = "/new/graph/pathway/"
    var url = baseURL+visType+name
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


function LoadPathway(pathwayId) { 
    serverAddr = NewGraphServer(name);
    wsURL = "ws://"+serverAddr;

   $("#cy").innerHTML = ""; 
    
    loadCy(); 
    

} 

var pwid = "hsa04630"
var pwid2 = "hsa04915"
var pwid3 = "hsa05200"

// current pathway
var logpwid = window.location.pathname.split("=")[1]; 

var numRuns = 0
function StartBenchmarks(){ 

    minSamples = 150

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

    

        
    /*
    suite.add("loadPathway "+pwid, function(deferred){
            name = "pathwaySelect="+pwid
            LoadPathway(name)
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


    suite.add("loadPathway "+pwid2, function(deferred){
            name = "pathwaySelect="+pwid2
            LoadPathway(name)
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

    suite.add("loadPathway "+pwid3, function(deferred){
            name = "pathwaySelect="+pwid3
            LoadPathway(name)
            defff = deferred
            numRuns = numRuns + 1 
            console.log(numRuns)
        
        },
        { 
            'defer': true,
            'minSamples': 80,
            'onComplete': function(e){
                console.log(e.currentTarget.name," is done running")
            }

        } 
    );

    */ 

    suite.add("set scale "+logpwid, function(deferred){
            setScale("log");       
            defff = deferred
        },
        {
            'minSamples': minSamples,
            'defer': true,
            'setup': function(){
                setScale("abs")
            },
    

        })

    suite.run({'async':true});
    benchmarked = true

} 



var defff

function deferAway(){ 
    
    try {
        defff.resolve()
        nodes = []; 
        edges = []; 

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
