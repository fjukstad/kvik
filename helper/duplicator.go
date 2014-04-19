package main

import (
    "log"
    "math/rand"
    "time"
    "strconv"
    "flag"  
    "os"
    . "os"
    "encoding/csv" 
    "fmt"
) 


func randomFloat() (float64){
    // seed 
    r := rand.New(rand.NewSource(time.Now().UnixNano())) 
    randomFloat := r.Float64()                              
    return randomFloat
}

func randomFloats(num int) [] string {
    var res [] string
    for i := 0; i < num; i++{
        f := strconv.FormatFloat(randomFloat()*1000, 'f', -1, 64)
        res = append(res, f) 
    }
    return res
} 


func main(){
    var numPairs = flag.Int("numpairs", 77, "number of cc pairs")
    var filename = flag.String("filename", "data.csv", "where to store data")
    var dsfile = flag.String("dsfile", "/Users/bjorn/stallo/data/exprs.csv",
                                "original dataset. where to get gene names")
    flag.Parse() 

    // Read header from original dataset
    exprsfile, err := os.Open(*dsfile)
    
    if err != nil {
        log.Panic(err)
    }

    defer exprsfile.Close() 

    reader := csv.NewReader(exprsfile) 
    header, err := reader.Read() 

    numGenes := len(header) - 1

    file, err := os.OpenFile(*filename, O_RDWR | O_CREATE, 0666)
    
    if err != nil { 
        log.Println("Could not open or create file", err)
        return 
    }


    log.Println("Writing dataset with ",numGenes,"genes and",*numPairs,
                    "cc pairs to",*filename)

    writer := csv.NewWriter(file) 

    // write header to file
    writer.Write(header) 
    
    var np int
    np = *numPairs
    for i := 0; i < *numPairs; i++ {
        p := strconv.Itoa(i)
        
        record := make([]string,0)
        record = append(record,p)
        record = append(record, randomFloats(numGenes)...)
        
        writer.Write(record) 

        p = p+"_1"
        record = make([]string,0)
        record = append(record,p)
        record = append(record, randomFloats(numGenes)...)
        writer.Write(record) 
        
        fmt.Print(i," of ",np," iterations done \r")
    }
    fmt.Print("\n")

    writer.Flush() 
    
    if err := writer.Error(); err != nil { 
        log.Println(err) 
    }
} 


// id, gene, gene, gene
// integer, gene, gene, gene
// integer_1, gene, gene, gene
// integer2, gene, gene, gene
// integer2_1, gene, gene, gene


