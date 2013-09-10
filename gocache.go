package main

import(
    "fmt"
    "net/http"
)

func main(){
    fmt.Println("hello"); 
}

func Get(url string) (resp *Response, err error) {
    

    

}

func getDirectory(url string) (dir string) {
	urlTokens := strings.Split(url, "/")
	strippedUrl := urlTokens[3:]
	dir = strings.Join(strippedUrl,"/")
	return 
}
