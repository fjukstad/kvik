package gocache

import (
    "testing"
    "os"
    //"log"
) 

// 
func Test_Get_1(t *testing.T) {
    clearCache()
    
    t.Log("On first request, url should not be found in cache")
    url := "http://blog.golang.org/error-handling-and-go"
    
    _, err := Get(url) 
    
    if err != nil {
        t.Error("Get failed: ", err)
        t.FailNow()
    }
    
}

func Test_Get_2(t *testing.T){
    url := "http://blog.golang.org/error-handling-and-go"
    _, err := Get(url)
    
    if err != nil{
        t.Error("Get Failed:",err)
        t.FailNow()
    }

}



func Test_Get_3(t *testing.T) {

    badurl := "htt://"

    _, err := Get(badurl)
    
    if err != nil {
        t.Log("Bad url got error:" , err)
        return
    }

    t.Error("Request with bad url succeeded") 

}

// Test benchmarking function
func Benchmark_Get_From_Empty_Cache(b *testing.B) {
    clearCache()
    url := "http://blog.golang.org/error-handling-and-go"
    for i := 0; i < b.N; i++{
        _, err := Get(url) 
        
        if err != nil{
            b.Error("Get failed in benchmark", err)
            b.FailNow()
        }
    }
}

func Benchmark_Get_From_Warmed_Up_Cache(b *testing.B){
    b.StopTimer() 
    url := "http://blog.golang.org/error-handling-and-go"
    
    _, _ = Get(url)

    b.StartTimer()
    
    for i := 0; i < b.N; i++{
        _, err := Get(url) 
        
        if err != nil{
            b.Error("Get failed in benchmark", err)
            b.FailNow()
        }
    }
    
}

func clearCache() {
    //log.Print("Clearing cache")
    os.RemoveAll("cache") 
}
