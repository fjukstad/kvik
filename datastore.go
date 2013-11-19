package main

import (
    "log"
    "flag"
)

func main() {
    
    var path = flag.String("path", "data" , "path where data files are stored")
    flag.Parse()
    
    log.Print("Generating dataset from directory: "+*path)
    
    newDataset(*path)

}

