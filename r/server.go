package r

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	keys     chan string
	sessions chan *Session
}

var cache map[string]string
var cachingEnabled bool

func (s *Server) Call(pkg, fun, args string) (string, error) {
	session := <-s.sessions
	res, err := session.Call(pkg, fun, args)
	s.sessions <- session
	return res, err
}

func (s Server) Get(key, format string) ([]byte, error) {
	session := <-s.sessions
	res, err := session.Get(key, format)
	s.sessions <- session
	return res, err
}

func (s *Server) CallHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Call read body failed", err)
		http.Error(w, "Call failed."+err.Error(), 500)
		return
	}

	call := Call{}
	err = json.Unmarshal(body, &call)
	if err != nil {
		fmt.Println("Call unmarshal failed", err)
		http.Error(w, "Call failed."+err.Error(), 500)
		return
	}

	log("Call:", call.Package, call.Function, call.Arguments)

	if cachingEnabled {
		key := call.cacheKey()
		res := cache[key]
		if res != "" {
			log("Cache hit")
			w.Write([]byte(res))
			return
		}
		log("Cache miss")
	}

	res, err := s.Call(call.Package, call.Function, call.Arguments)
	if err != nil {
		fmt.Println("Call failed", err)
		http.Error(w, "Call failed."+err.Error(), 500)

	}
	w.Write([]byte(res))

	if cachingEnabled {
		key := call.cacheKey()
		cache[key] = res
	}

	return
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]
	format := vars["format"]

	log("Get", key, format)

	res, err := s.Get(key, format)
	if err != nil {
		fmt.Println("Get failed", err)
		http.Error(w, "Get failed."+err.Error(), 500)
		return
	}

	w.Write(res)
}

func (s *Server) EnableCaching() {
	cachingEnabled = true
	cache = make(map[string]string, 0)
}

func (s *Server) Start(port string) error {

	router := mux.NewRouter()
	router.HandleFunc("/call", s.CallHandler)
	router.HandleFunc("/get/{key}/{format}", s.GetHandler)
	http.Handle("/", router)

	return http.ListenAndServe(port, router)

}
