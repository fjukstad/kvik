package gocache

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
    "path"
    "io"
)

type Entry struct{
    Response *http.Response
    Content string
}


func main(){
    fmt.Println("Testing web cache");
    url := "http://blog.golang.org/error-handling-and-go"
    Get(url)
}

func Get(url string) (resp *http.Response, err error) {
    
    resp, err = getFromCache(url) 
   
    if err != nil{
        log.Println(url,"could not be served out of cache.", err)
        return getFromWeb(url) 
    }
    
    log.Print("Request to ", url, " served out of cache") 
    return 
}

// Fetch contents of url from web and write to cache
func getFromWeb(url string) (resp *http.Response, err error){
    resp, err = http.Get(url) 
    
    if err != nil{
        return
    }
    var resp2 http.Response
    resp2 = *resp

    // Draining the body
    
    body, err := ioutil.ReadAll(resp.Body) 

    if err != nil {
        log.Print("Reading response body went bad. ", err); 
    
    }

    writeToCache(url, body, &resp2)   

    resp.Body = nopCloser{bytes.NewBufferString(string(body))} 


    return 
}

func writeToCache(url string, body []byte, resp *http.Response){
    
    cacheEntry := generateCacheEntry(resp, body)

    b, err := json.Marshal(cacheEntry)

    if err != nil {
        log.Println("Could not marshal response ", err)
        return 
    }
    

    filename := getFilePath(url)
    
    // Set up any directory needed to write the file. 
    // e.g. vg.no/first/second/file will be stored as
    // cache/first/second/file 

    err = createDirectories(filename)
    if err != nil{
        pe, _ := err.(*os.PathError) 

        if ! strings.Contains(pe.Error(),"file exists") {
            log.Println("Could not create directories ", err)
            return
        }
    }

    // Create file 
    file, err := os.Create(filename)
    if err != nil {
        log.Println("could not create cache file ", err)
        return 
    }

    // Close and check for error on exit 
    defer func() {
        if err := file.Close(); err != nil {
            
            log.Println("Could not close file ", err)
        }
    }()

    // Write file 
    _, err = file.Write(b)
    if err != nil{
        log.Println("Could not write to cache file ", err)
    }
    
    log.Println(url, "successfully written to cache"); 
}

func generateCacheEntry(resp *http.Response, body []byte ) Entry {

    /*
    body, err := ioutil.ReadAll(resp.Body) 

    if err != nil {
        log.Print("Reading response body went bad. ", err); 
    
    }
    */

    Response := resp
    Content := string(body)
    entry := Entry{Response, Content} 
    


    // Cannot marshal the body from get go
    entry.Response.Body = nil

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
    
    entry, err := readFromFile(file)
    
    if err != nil{
        return
    }

    resp = entry.GenerateHttpResponse() 

    return

} 

// Read cache entry from file
func readFromFile(file *os.File) (entry *Entry, err error) {
    fileinfo, err := file.Stat() 
    var size int
    size = int(fileinfo.Size())

    buf := make([]byte, size)
    
    n, err := file.Read(buf)

    if err != nil{
        log.Println("Reading file got:", err, "after",n,"bytes"); 
        return entry, errors.New("Reading file went wrong") 
    }
        

    // Must trim buffer before unmarshaling it. This is because of 
    // the unmarshaling failing if entire buffer is returned
    buf = bytes.Trim(buf[0:], string(0))    


    entry = new(Entry)
    err = json.Unmarshal(buf, entry)

    if err != nil {
        log.Print("Unmarshaling gone wrong: ", err) 
        return entry, errors.New("Unmarshaling gone wrong")
    }

    return
}

func (entry *Entry) Print () {
    
    //n := bytes.Index(entry.Content, []byte{0})
    //content := string(entry.Content[:n]) 
    log.Print("Content: ", entry.Content) 

}

func (entry *Entry) GenerateHttpResponse() (resp *http.Response){

    resp = entry.Response
    resp.Body = nopCloser{bytes.NewBufferString(entry.Content)} 


    return resp

}


func getFilePath(url string) (dir string) {
	urlTokens := strings.Split(url, "/")
    // 2 because we need to strip away 'http:', ' '
    strippedUrl := urlTokens[2:] 
	dir = "cache/"+strings.Join(strippedUrl,"/")
	return 
}

func createDirectories(filename string) error {
    
    filepath := path.Dir(filename)
    directories := strings.Split(filepath, "/") 
    
    p := "" 
    for i := range directories {
        p += directories[i] + "/"
        err := os.Mkdir(p, 0755) 

        if err != nil {
            pe, _ := err.(*os.PathError) 
            
            // if folder exists, continue to the next one
            if ! strings.Contains(pe.Error(),"file exists") {
                log.Println("Mkdir failed miserably: ", err) 
                return err
            }
        }
    }

    return nil
}


/* Below are from: 
    https://groups.google.com/forum/#!topic/golang-nuts/J-Y4LtdGNSw
*/

type nopCloser struct { 
    io.Reader 
} 

func (nopCloser) Close() error { 
    return nil
} 
