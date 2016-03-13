package r

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var r *rand.Rand
var rootDir string

type Session struct {
	Key     string
	Dir     string
	Output  string
	Command string
}

type RCall struct {
	Package   string
	Function  string
	Arguments string
}

func Init(dir, packages string) error {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	rootDir = dir

	if packages != "" {
		err := installPackages(packages)
		if err != nil {
			fmt.Println("could not install packages")
			return err
		}
	}

	pkgs, err := InstalledPackages()
	if err != nil {
		fmt.Println("Could not get installed packages")
		return err
	}

	err = ioutil.WriteFile(dir+"/r-packages.json", pkgs, 0755)
	if err != nil {
		fmt.Println("Could not write r packages file")
		return err
	}
	return nil
}

// Installs the given packages on the server. packages is a space
// separated list with the packages.
func installPackages(packages string) error {
	cmd := exec.Command("R", "-q", "-e", "install.packages(c("+packages+"),  repos='http://cran.us.r-project.org')")

	var out bytes.Buffer
	cmd.Stdout = &out

	fmt.Println(out.String())
	err := cmd.Run()
	if err != nil {
		fmt.Println("Could not install packages:", out.String())
		return err
	}
	return nil
}

// Installs the package found at src. Should be a tgz ball.
func InstallPackageFromSource(src string) (string, error) {

	cmd := exec.Command("R", "CMD", "INSTALL", src)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Could not install local package:", out.String())
		return "", err
	}
	return out.String(), nil

}

// Executes function. Returns tmp key for use in Get
func Call(pkg, fun, args string) (*Session, error) {

	h := md5.New()
	k := h.Sum([]byte(strconv.Itoa(r.Int())))

	hash := hex.EncodeToString(k)
	key := ".s" + hash

	wd := rootDir + "/" + key
	err := os.MkdirAll(wd, 0755)
	if err != nil {
		return nil, err
	}

	err = os.Chdir(wd)

	if err != nil {
		return nil, err
	}

	// Replace argument names with real argument names using the
	// keys from previous calls. Also load data from previous
	// calls before running cmd

	argList := strings.Split(args, ",")
	loadArgs := []string{}

	if argList[0] != "" {
		finalArgs := []string{}

		for _, arg := range argList {

			// special case when we use a column vector as argument value
			// e,g. d = c(2,1,3,4) only d=c(2, is added. we need to
			// append the rest to get the full vector.
			if len(strings.Split(arg, "=")) == 1 {
				finalArgs[len(finalArgs)-1] = finalArgs[len(finalArgs)-1] + "," + arg
				continue
			}

			argName := strings.Split(arg, "=")[0]
			argVal := strings.Split(arg, "=")[1]
			if strings.HasPrefix(argVal, ".s") {
				loadArgs = append(loadArgs, "load('"+rootDir+"/"+argVal+"/.RData');")

				argVal = strings.TrimPrefix(argVal, ".")

			}
			finalArgs = append(finalArgs, argName+"="+argVal)
		}

		args = strings.Join(finalArgs, ",")
	}

	varName := strings.TrimPrefix(key, ".")
	command := varName + "=" + pkg + "::" + fun + "(" + args + "); " + varName

	if len(loadArgs) > 0 {
		loadString := strings.Join(loadArgs, "")
		command = loadString + command
	}

	cmd := exec.Command("R", "--save", "-q", "-e", command)
	cmd.Dir = wd

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()

	return &Session{key, wd, out.String(), command}, err
}

// get results for a func call
func Get(key, format string) ([]byte, error) {

	dir := rootDir + "/" + key
	err := os.Chdir(dir)
	if err != nil {
		return nil, err
	}

	varName := strings.TrimPrefix(key, ".")

	extension := "." + format
	_, err = os.Stat("output" + extension)
	if err == nil {
		return ioutil.ReadFile(dir + "/output" + extension)
	}

	var command string
	if format == "csv" {
		command = "write.table(" + varName + ", sep=',', row.names=FALSE, file='output" + extension + "')"
	} else if format == "json" {
		command = "js=jsonlite::toJSON(" + varName + "); write(js, file='output" + extension + "')"
	} else if format == "pdf" {
		return ioutil.ReadFile(dir + "/Rplots.pdf")
	} else if format == "png" {
		cmd := exec.Command("pdftoppm", "-png", dir+"/Rplots.pdf", "plot")

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			fmt.Println("Could not convert rplot to png:", out.String())
			return nil, err
		}

		return ioutil.ReadFile(dir + "/plot-1.png")
	} else {
		return nil, errors.New("Unknown format")
	}

	cmd := exec.Command("R", "--save", "-q", "-e", command)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		fmt.Println("ERROR", out.String())
		return nil, err
	}

	return ioutil.ReadFile(dir + "/output" + extension)

}

func InstalledPackages() ([]byte, error) {

	pkg := "utils"
	fun := "installed.packages"
	args := ""

	s, err := Call(pkg, fun, args)
	if err != nil {
		fmt.Println("could not get installed packages")
		return nil, err
	}

	pkg = "base"
	fun = "as.data.frame"
	args = "x=" + s.Key

	s, err = Call(pkg, fun, args)
	if err != nil {
		fmt.Println("could not get installed packages")
		return nil, err
	}

	return Get(s.Key, "json")
}

type PackageInfo struct {
	Package               string
	LibPath               string
	Version               string
	Priority              string
	Depends               string
	Imports               string
	LinkingTo             string
	Suggests              string
	Enhances              string
	Licence               string
	License_is_FOSS       string
	Licence_restricts_use string
	OS_type               string
	MD5sum                string
	NeedsCompilation      string
	Built                 string
}

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

// Uploads and installs package to remote R Server
func (s *Server) UploadPackage(src string) error {
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

	url := "http://" + s.Addr + "/install"

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

	http.Handle("/", r)

	printTime()
	fmt.Println("Starting server at ", port)

	return http.ListenAndServe(port, r)

}
