package main

import(
    "fmt"
    "net/http"
)

func main(){
    fmt.Println("hello"); 
}

func Get(url string) (resp *Response, err error) {
    
    resp, err := getFromCache(url) 
    
    if err != nil{
        return http.Get(url) 
    }

    return   
}

func getFromCache(url string) (resp *Response, err error) {
    

    
} 


func getDirectory(url string) (dir string) {
	urlTokens := strings.Split(url, "/")
    // 3 because we need to strip away 'http:', ' ', and hostname
    strippedUrl := urlTokens[3:] 
	dir = strings.Join(strippedUrl,"/")
	return 
}
