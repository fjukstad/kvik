package main

import(
    "fmt"
    "net/http"
    "errors"
    "os"
    "strings"
    "log"
    "encoding/json" 
    "bytes"
    "io/ioutil"
)

type Entry struct{
    Response *http.Response
    Content []byte
}


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
    
    cacheEntry := generateCacheEntry(resp)

    b, err := json.Marshal(cacheEntry)

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

func generateCacheEntry(resp *http.Response) Entry {

    body, err := ioutil.ReadAll(resp.Body) 

    if err != nil {
        log.Print("Reading response body went bad. ", err); 
    }
    
    //n := bytes.Index(body, []byte{0})
    //content := string(body[:n]) 
    
    //entry.Content = body
    //entry.Response = res

    entry := Entry{resp, body}
    
    return entry

}


// Try to fetch contents of url from cache
func getFromCache(url string) (resp *http.Response, err error) {
    
    filename := getFilePath(url)
    
    file, err := os.Open(filename) 
    if err != nil{
        err = errors.New("File '"+filename+"' not Found") 
        return 
    }
    
    entry := readFromFile(file)

    fmt.Println(file)
    log.Println(entry) 
    return

} 

// Read cache entry from file
func readFromFile(file *os.File) (entry *Entry) {
    size := 8192*32

    buf := make([]byte, size)
    
    n, err := file.Read(buf)

    if err != nil{
        log.Println("Reading file got:", err, "after",n,"bytes"); 
        return
    }
        
    log.Println("Read",n,"bytes")

    log.Println(buf); 


    // Must trim buffer before unmarshaling it. This is because of 
    // the unmarshaling failing if entire buffer is returned
    buf = bytes.Trim(buf[0:], string(0))    


    log.Println(buf); 

    entry = new(Entry)

    err = json.Unmarshal(buf, entry)

    if err != nil {
        log.Print("Unmarshaling gone wrong: ", err) 
       // log.Print("buf is now: ", buf)
       // log.Print("and resp is :", entry); 
    }

    entry.Print()

    return
}

func (entry *Entry) Print () {
    
    n := bytes.Index(entry.Content, []byte{0})
    log.Print("God damn n: ",n)
    content := string(entry.Content) 
    
    log.Print("Content: ", content) 

    log.Panic("stop here")



}


func getFilePath(url string) (dir string) {
	urlTokens := strings.Split(url, "/")
    // 3 because we need to strip away 'http:', ' ', and hostname
    strippedUrl := urlTokens[3:] 
	dir = strings.Join(strippedUrl,"/")
	return 
}
