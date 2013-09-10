package main

import(
    "fmt"
    "net/http"
    "errors"
    "os"
    "strings"
    "log"
    "encoding/json" 
)

func main(){
    fmt.Println("hello");
    url := "http://blog.golang.org/error-handling-and-go"
    Get(url)
}

func Get(url string) (resp *http.Response, err error) {
    
    resp, err = getFromCache(url) 
    
    if err != nil{
        log.Println(url,"could not be served out of cache.", err)
        return getFromWeb(url) 
    }

    return   
}

// Fetch contents of url from web and write to cache
func getFromWeb(url string) (resp *http.Response, err error){
    resp, err = http.Get(url) 
    
    if err != nil{
        return
    }
    writeToCache(url, resp)   
    return 
}

func writeToCache(url string, resp *http.Response){
    
    log.Println("Writing response to cache")

    b, err := json.Marshal(resp)

    if err != nil {
        log.Println(err)
    }
    
    filename := getFilePath(url)
    file, err := os.Create(filename)
    
    if err != nil {
        log.Println(err)
        return
    }

    // Close and check for error on exit 
    defer func() {
        if err := file.Close(); err != nil {
            log.Println(err)
        }
    }()
        
    _, err = file.Write(b)
    if err != nil{
        log.Println(err)
    }
    
    log.Println("Successfully written to cache"); 
}

// Try to fetch contents of url from cache
func getFromCache(url string) (resp *http.Response, err error) {
    
    filename := getFilePath(url)
    
    file, err := os.Open(filename) 
    if err != nil{
        err = errors.New("File '"+filename+"' not Found") 
        return 
    }
    
    resp = readFromFile(file)

    fmt.Println(file)
    return

} 

// Read cache entry from file
func readFromFile(file *os.File) (resp *http.Response) {
    
    buf := make([]byte, 1024)
    
    for {
        _, err := file.Read(buf)
        if err != nil{
            break
        }
    }

    err := json.Unmarshal(buf,resp)

    if err != nil {
        log.Print(err) 
    }


    return
}

func getFilePath(url string) (dir string) {
	urlTokens := strings.Split(url, "/")
    // 3 because we need to strip away 'http:', ' ', and hostname
    strippedUrl := urlTokens[3:] 
	dir = strings.Join(strippedUrl,"/")
	return 
}
