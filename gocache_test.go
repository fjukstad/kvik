package gocache

import (
    "testing"
    "os"
    "log"
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


func clearCache() {
    log.Print("Clearing cache")
    os.RemoveAll("cache") 
}
