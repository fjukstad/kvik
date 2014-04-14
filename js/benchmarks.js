/*
    var suite = new Benchmark.Suite('test', { 
        'onComplete': function(){
            for(i = 0; i < suite.length; i++){
                bench = suite[i];
                console.log(bench)
                PostResults(bench.name, JSON.stringify(bench.stats))
            }

        }
    });
    suite.add("init graph",LongRunningTask);
    suite.add('init graph nocache', function(){
                    ClearBackendCache();
                    LongRunningTask();
                    }
                ,{
                'maxTime':25,
                //'minSamples': 10
            });
    
    //for(var i = 0; i < 5; i++){
        suite.run({'async':false}) 
    //}

    suite.run();

    console.profileEnd()

    console.log("done")

    console.log(console.profiles())
*/
function PostResults(method,data) {
    var url = "http://localhost:4444/result/"+method
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
    console.log("loading pathway:",name)
    serverAddr = NewGraphServer(name);
    wsURL = "ws://"+serverAddr;

   $("#cy").innerHTML = ""; 
    
    loadCy(); 
    
    console.log("hepp")

} 


function StartBenchmarks(){ 

    var suite = new Benchmark.Suite('test', { 
        'onComplete': function(){
            for(i = 0; i < suite.length; i++){
                bench = suite[i];
                console.log(bench)
                PostResults(bench.name, JSON.stringify(bench.stats))
            }

        }
    });

    suite.add("hsa04630 with vis", function(deferred){
            name = "pathwaySelect=hsa04630"
            LoadPathway(name)
            defff = deferred
        },
        { 
            'defer': true
        } 
        );

    suite.run({'async':true});

    benchmarked = true

        console.log(suite) 
    
} 

var defff

function deferAway(){ 
    
    try {
        defff.resolve()
    }
    catch(e){
        console.log(e)
    } 
} 
