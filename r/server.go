package r

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Addr     string
	Username string
	Server   string
}

// Remote call.
func (s *Server) Call(pkg, fun, args string) (string, error) {
	url := "http://" + s.Addr + "/call"

	c := RCall{pkg, fun, args}

	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	postArgs := strings.NewReader(string(b))

	resp, err := http.Post(url, "application/json", postArgs)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if strings.Contains(string(body), "exit status 1") {
		return "exit status 1", errors.New(string(body))
	}

	return string(body), nil
}

// Remote get results
func (s *Server) Get(key, format string) ([]byte, error) {

	resp, err := http.Get("http://" + s.Addr + "/get/" + key + "/" + format)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil

}

// Uploads something to a  remote R Server. Something can be e.g. a package to
// install or script to execute.
func (s *Server) Upload(src, path string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	file.Close()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return err
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		return err
	}

	url := "http://" + s.Addr + "/" + path

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Response:", string(responseBody))
	return nil
}

func (s *Server) InstalledPackages() ([]byte, error) {
	pkg := "utils"
	fun := "installed.packages"
	args := ""

	session, err := s.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("could not get installed packages")
		return nil, err
	}

	pkg = "base"
	fun = "as.data.frame"
	args = "x=" + session

	session, err = s.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("could not get installed packages")
		return nil, err
	}

	return s.Get(session, "json")

}

var cache map[string]string
var logTime = "02-01-2006 15:04:05.000"

func CallHandler(w http.ResponseWriter, r *http.Request) {
	printTime()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Call failed %s", err)
		fmt.Println(err)
		return
	}

	call := RCall{}
	err = json.Unmarshal(body, &call)
	if err != nil {
		fmt.Fprintf(w, "Call failed. %s", body)
		fmt.Println(err)
		return
	}

	fmt.Println("Call:", call.Package, call.Function, call.Arguments)

	cacheKey := call.Package + call.Function + call.Arguments

	if cache[cacheKey] != "" {
		printTime()
		fmt.Println("CACHE HIT")
		w.Write([]byte(cache[cacheKey]))
		return
	} else {
		printTime()
		fmt.Println("CACHE MISS")
	}

	s, err := Call(call.Package, call.Function, call.Arguments)
	if err != nil {
		fmt.Println("ERROR", err)
		fmt.Println(s.Output)
		fmt.Fprintf(w, "Call failed. %s returned %s", s.Output, err)
		return
	}

	if cache == nil {
		cache = make(map[string]string, 0)
	}

	cache[cacheKey] = s.Key

	w.Write([]byte(s.Key))

}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	printTime()
	vars := mux.Vars(r)
	key := vars["key"]
	format := vars["format"]

	fmt.Println("Get: key:", key, "format", format)

	res, err := Get(key, format)

	if err != nil {
		fmt.Println(err)
		fmt.Println(string(res))
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nothing to see here. Try /call, /get/{key}/{format} or /install"))
}

func InstallHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fileLoc := "/tmp/file"

	out, err := os.Create(fileLoc)
	if err != nil {
		fmt.Fprintf(w, "Failed to open the file for writing")
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	printTime()
	fmt.Println("Installing package", header.Filename)

	output, err := InstallPackageFromSource(fileLoc)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Fprintf(w, "%s", output)

}

func printTime() {
	fmt.Print(time.Now().Format(logTime), " ")
}

func StartServer(port, tmpDir string) error {
	err := Init("/tmp/kvik", "")
	if err != nil {
		fmt.Println(err)
		return err
	}

	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/call", CallHandler)
	r.HandleFunc("/get/{key}/{format}", GetHandler)
	r.HandleFunc("/install", InstallHandler)
	r.HandleFunc("/script", ScriptHandler)

	http.Handle("/", r)

	printTime()
	fmt.Println("Starting server at ", port)

	return http.ListenAndServe(port, r)

}
