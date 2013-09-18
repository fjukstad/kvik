# gocache
A HTTP cache written in go. 


**Note**: When webpages are cached they are not updated unless they are manually removed from the filesystem. 
Cached webpages are stored in a _cache_ directory in the current working directory of your go app. 

# Installing
  
    go get github.com/fjukstad/gocache
  
# Using it
Add
  
    github.com/fjukstad/gocache
  
to your imports and replace your
  
    http.Get(url)

with


    gocache.Get(url) 
  
and you're good to go. The files are stored in the working directory of your go-app under a _cache_ directory. 


# Example 

```go
package main
    
import (
        "github.com/fjukstad/gocache"
        )
  
func main() {

    url := "http://blog.golang.org/error-handling-and-go"
    resp := gocache.Get(url) 
    
    // do whatever you'd like with resp 
}
    
