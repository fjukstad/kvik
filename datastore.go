package main

import (
    "log"
    "flag"
/*
    "runtime/pprof"
    "os"
*/
)

func main() {
    
    var path = flag.String("path", "data" , "path where data files are stored")
    flag.Parse()
    
    log.Print("Generating dataset from directory: "+*path)
    

    newDataset(*path)
/*
    
// Profiling
// $ go tool pprof /home/bfj001/master/src/bin/datastore memprofile.prof 
// (pprof) top5
    f, err := os.Create("memprofile.prof")
    if err != nil {
        log.Fatal(err)
    }
    pprof.WriteHeapProfile(f)
    f.Close()
    return
*/

}

