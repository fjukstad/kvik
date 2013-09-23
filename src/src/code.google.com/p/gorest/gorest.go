//Copyright 2011 Siyabonga Dlamini (siyabonga.dlamini@gmail.com). All rights reserved.
//
//Redistribution and use in source and binary forms, with or without
//modification, are permitted provided that the following conditions
//are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer
//     in the documentation and/or other materials provided with the
//     distribution.
//
//THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
//IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
//OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
//IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
//PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS;
//OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
//WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
//OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
//ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gorest

import (
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
)

type GoRestService interface {
	ResponseBuilder() *ResponseBuilder
}

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

type endPointStruct struct {
	name                 string
	requestMethod        string
	signiture            string
	muxRoot              string
	root                 string
	nonParamPathPart     map[int]string
	params               []param //path parameter name and position
	queryParams          []param
	signitureLen         int
	paramLen             int
	inputMime            string
	outputType           string
	outputTypeIsArray    bool
	outputTypeIsMap      bool
	postdataType         string
	postdataTypeIsArray  bool
	postdataTypeIsMap    bool
	isVariableLength     bool
	parentTypeName       string
	methodNumberInParent int
	role                 string
}

type restStatus struct {
	httpCode int
	reason   string //Especially for code in range 4XX to 5XX
}

func (err restStatus) String() string {
	return err.reason
}

type serviceMetaData struct {
	template     interface{}
	consumesMime string
	producesMime string
	root         string
	realm        string
}

var restManager *manager
var handlerInitialised bool

type manager struct {
	serviceTypes map[string]serviceMetaData
	endpoints    map[string]endPointStruct
}

func newManager() *manager {
	man := new(manager)
	man.serviceTypes = make(map[string]serviceMetaData, 0)
	man.endpoints = make(map[string]endPointStruct, 0)
	return man
}
func init() {
	RegisterMarshaller(Application_Json, NewJSONMarshaller())

}

//Registeres a service on the rootpath.
//See example below:
//
//	package main
//	import (
// 	   "code.google.com/p/gorest"
//	        "http"
//	)
//	func main() {
//	    gorest.RegisterService(new(HelloService)) //Register our service
//	    http.Handle("/",gorest.Handle())
// 	   http.ListenAndServe(":8787",nil)
//	}
//
//	//Service Definition
//	type HelloService struct {
//	    gorest.RestService `root:"/tutorial/"`
//	    helloWorld  gorest.EndPoint `method:"GET" path:"/hello-world/" output:"string"`
//	    sayHello    gorest.EndPoint `method:"GET" path:"/hello/{name:string}" output:"string"`
//	}
//	func(serv HelloService) HelloWorld() string{
// 	   return "Hello World"
//	}
//	func(serv HelloService) SayHello(name string) string{
//	    return "Hello " + name
//	}
func RegisterService(h interface{}) {
	RegisterServiceOnPath("", h)
}

//Registeres a service under the specified path.
//See example below:
//
//	package main
//	import (
//	    "code.google.com/p/gorest"
//	        "http"
//	)
//	func main() {
//	    gorest.RegisterServiceOnPath("/rest/",new(HelloService)) //Register our service
//	    http.Handle("/",gorest.Handle())
//	    http.ListenAndServe(":8787",nil)
//	}
//
//	//Service Definition
//	type HelloService struct {
//	    gorest.RestService `root:"/tutorial/"`
//	    helloWorld  gorest.EndPoint `method:"GET" path:"/hello-world/" output:"string"`
//	    sayHello    gorest.EndPoint `method:"GET" path:"/hello/{name:string}" output:"string"`
//	}
//	func(serv HelloService) HelloWorld() string{
//	    return "Hello World"
//	}
//	func(serv HelloService) SayHello(name string) string{
//	    return "Hello " + name
//	}
func RegisterServiceOnPath(root string, h interface{}) {
	//We only initialise the handler management once we know gorest is being used to hanlde request as well, not just client.
	if !handlerInitialised {
		restManager = newManager()
		handlerInitialised = true
	}

	if root == "/" {
		root = ""
	}

	if root != "" {
		root = strings.Trim(root, "/")
		root = "/" + root
	}

	registerService(root, h)
}

//ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
func (man *manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	url_, err := url.QueryUnescape(r.URL.RequestURI())

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Internal Server Error: Could not serve page: ", r.Method, url_)
			log.Println(rec)
			log.Printf("%s", debug.Stack())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	if err != nil {
		log.Println("Could not serve page: ", r.Method, r.URL.RequestURI(), "Error:", err)
		w.WriteHeader(400)
		w.Write([]byte("Client sent bad request."))
		return
	}

	if ep, args, queryArgs, xsrft, found := getEndPointByUrl(r.Method, url_); found {

		ctx := new(Context)
		ctx.writer = w
		ctx.request = r
		ctx.args = args
		ctx.queryArgs = queryArgs
		ctx.xsrftoken = xsrft

		data, state := prepareServe(ctx, ep)

		if state.httpCode == http.StatusOK {
			switch ep.requestMethod {
			case POST, PUT, DELETE, HEAD, OPTIONS:
				{
					if ctx.responseCode == 0 {
						w.WriteHeader(getDefaultResponseCode(ep.requestMethod))
					} else {
						if !ctx.dataHasBeenWritten {
							w.WriteHeader(ctx.responseCode)
						}
					}
				}
			case GET:
				{
					if ctx.responseCode == 0 {
						if !ctx.responseMimeSet {
							w.Header().Set("Content-Type", man.getType(ep.parentTypeName).producesMime)
						}
						w.WriteHeader(getDefaultResponseCode(ep.requestMethod))
					} else {
						if !ctx.dataHasBeenWritten {
							if !ctx.responseMimeSet {
								w.Header().Set("Content-Type", man.getType(ep.parentTypeName).producesMime)
							}
							w.WriteHeader(ctx.responseCode)
						}
					}

					if !ctx.overide {
						w.Write(data)
					}

				}
			}

		} else {
			log.Println("Problem with request. Error:", r.Method, state.httpCode, state.reason, "; Request: ", r.URL.RequestURI())
			w.WriteHeader(state.httpCode)
			w.Write([]byte(state.reason))
		}

	} else {
		log.Println("Could not serve page, path not found: ", r.Method, url_)
		//		println("Could not serve page, path not found: ", r.Method, url_)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("The resource in the requested path could not be found."))
	}

}

func (man *manager) getType(name string) serviceMetaData {

	return man.serviceTypes[name]
}
func (man *manager) addType(name string, i serviceMetaData) string {
	for str, _ := range man.serviceTypes {
		if name == str {
			return str
		}
	}

	man.serviceTypes[name] = i
	return name
}
func (man *manager) addEndPoint(ep endPointStruct) {
	man.endpoints[ep.requestMethod+":"+ep.signiture] = ep
}

//Registeres the function to be used for handling all requests directed to gorest.
func HandleFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving URL : ", r.Method, r.URL.RequestURI())
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Internal Server Error: Could not serve page: ", r.Method, r.RequestURI)
			log.Println(rec)
			log.Printf("%s", debug.Stack())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	restManager.ServeHTTP(w, r)
}

//Runs the default "net/http" DefaultServeMux on the specified port.
//All requests are handled using gorest.HandleFunc()
func ServeStandAlone(port int) {
	http.HandleFunc("/", HandleFunc)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func _manager() *manager {
	return restManager
}
func Handle() *manager {
	return restManager
}

func getDefaultResponseCode(method string) int {
	switch method {
	case GET, PUT, DELETE:
		{
			return 200
		}
	case POST:
		{
			return 202
		}
	default:
		{
			return 200
		}
	}

	return 200
}
