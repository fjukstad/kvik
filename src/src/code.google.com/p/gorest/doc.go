//
//A RESTful style web-services framework for the Go language. </br>
//Creating services in Go is straight forward, GoRest? takes this a step further by adding a layer that 
//makes tedious tasks much more automated and avoids regular pitfalls. <br/> 
//This gives you the opportunity to focus more on the task at hand... minor low-level http handling.<br/>
//
//
//Example usage below:
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
package gorest
