package r

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	keys     chan string
	sessions chan *Session
}

var cache map[string]string
var cachingEnabled bool
var cacheMutex *sync.Mutex

func (s *Server) Call(pkg, fun, args string) (string, error) {
	session := <-s.sessions
	res, err := session.Call(pkg, fun, args)
	if err != nil {
		var new_err error
		session.cmd.Process.Kill()
		session, new_err = NewSession(session.id)
		if new_err != nil {
			return "", errors.Wrap(err, "Could not start new R Session")
		}
	}
	s.sessions <- session
	return res, err
}

func (s Server) Get(key, format string) ([]byte, error) {
	session := <-s.sessions
	res, err := session.Get(key, format)
	if err != nil {
		var new_err error
		session.cmd.Process.Kill()
		session, new_err = NewSession(session.id)
		if new_err != nil {
			return []byte{}, errors.Wrap(err, "Could not start new R Session")
		}
	}
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
		cacheMutex.Lock()
		res := cache[key]
		cacheMutex.Unlock()

		if res != "" {
			log("Cache hit")
			w.Write([]byte(res))
			return
		}
		log("Cache miss")
	} else {
		log("No cache")
	}

	res, err := s.Call(call.Package, call.Function, call.Arguments)
	if err != nil {
		http.Error(w, "Call failed: "+err.Error(), 500)
		return
	}
	w.Write([]byte(res))
	if cachingEnabled {
		key := call.cacheKey()
		cacheMutex.Lock()
		cache[key] = res
		cacheMutex.Unlock()
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
	cacheMutex = &sync.Mutex{}
}

func (s *Server) Start(port string) error {

	router := mux.NewRouter()
	router.HandleFunc("/call", s.CallHandler)
	router.HandleFunc("/get/{key}/{format}", s.GetHandler)
	http.Handle("/", router)

	return http.ListenAndServe(port, router)

}
