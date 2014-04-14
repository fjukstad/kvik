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


