package r

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var rootDir string
var timeout int64 = 600000

type Session struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	id     int
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (rs *Session) Call(pkg, fun, args string) (string, error) {

	key := "." + randSeq(5)
	wd := rootDir + "/" + key
	varName := strings.TrimPrefix(key, ".")

	err := os.MkdirAll(wd, 0755)
	if err != nil {
		return "", errors.Wrap(err, "Could not create tmp dir for results")
	}

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

			if strings.HasPrefix(argVal, ".") {
				loadArgs = append(loadArgs, "load('"+rootDir+"/"+argVal+"/.RData');")
				argVal = strings.TrimPrefix(argVal, ".")
			}
			finalArgs = append(finalArgs, argName+"="+argVal)
		}
		args = strings.Join(finalArgs, ",")
	}

	command := varName + "=" + pkg + "::" + fun + "(" + args + ");"

	if len(loadArgs) > 0 {
		loadString := strings.Join(loadArgs, "")
		command = loadString + command
	}

	// load desired package before calling function (makes sure
	// that we load packages it depends on)
	command = "library(" + pkg + ");" + command
	command = "rm(list=ls());" + "setwd(\"" + wd + "\");pdf();" + command + "save.image();dev.off();\n"

	_, err = io.WriteString(rs.stdin, command)
	if err != nil {
		return "", err
	}

	res, err := rs.getResult(key, wd+"/.RData")
	if err != nil {
		return "", err
	}

	return res, err
}

func (rs *Session) getResult(key, filename string) (string, error) {
	keys := make(chan string, 1)
	errCh := make(chan string, 1)

	go func(ch chan string, filename, key string) {
		for {
			finfo, err := os.Stat(filename)
			if err == nil && finfo.Size() > 0 {
				ch <- key
				return
			}
		}
	}(keys, filename, key)

	go func(ch chan string, stderr io.ReadCloser) {
		buf := new(bytes.Buffer)
		for {
			buf.ReadFrom(stderr)
			errmsg := buf.String()
			if errmsg != "" {
				ch <- errmsg
				return
			}
		}
	}(errCh, rs.stderr)

	select {
	// call successful no errors
	case k := <-keys:
		return k, nil
	case o := <-errCh:
		fmt.Println(o)
		return "", errors.New(o)
	case <-time.After(time.Duration(timeout) * time.Second):
		return "", errors.New("R Call timeout.")
	}

	return key, nil
}

func (rs *Session) Get(key, format string) ([]byte, error) {
	wd := rootDir + "/" + key

	cmd := "rm(list=ls());" + "setwd(\"" + wd + "\");" + "load(\".RData\");\n"

	_, err := io.WriteString(rs.stdin, cmd)
	if err != nil {
		fmt.Println("Could not write to R process", err)
		return []byte{}, err
	}

	varName := strings.TrimPrefix(key, ".")

	filename := wd + "/output." + format

	var command string
	if format == "json" {
		command = "js=jsonlite::toJSON(" + varName + "); write(js, file='" + filename + "');\n"
	} else if format == "csv" {
		command = "write.table(" + varName + ", sep=',', row.names=FALSE, file='" + filename + "');\n"
	} else if format == "pdf" {
		file, err := ioutil.ReadFile(wd + "/Rplots.pdf")
		if err != nil {
			return nil, errors.Wrap(err, "Could not read pdf file")
		}
		// if file contains magic end return it, if not wait
		// for R to complete writing the file.
		for !strings.Contains(string(file), "%%EOF") {
			time.Sleep(500 * time.Millisecond)
			file, err = ioutil.ReadFile(wd + "/Rplots.pdf")
			if err != nil {
				return nil, errors.Wrap(err, "Could not read pdf file:")
			}
		}
		return file, err
	} else if format == "png" {
		cmd := exec.Command("pdftoppm", "-png", wd+"/Rplots.pdf", wd+"/plot")
		//cmd := exec.Command("convert", wd+"/Rplots.pdf", wd+"/plot-1.png")

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return nil, errors.Wrap(err, "Could not convert Rplot pdf to png")
		}
		return ioutil.ReadFile(wd + "/plot-1.png")
	} else {
		return nil, errors.Wrap(err, "Unknown format:")
	}

	_, err = io.WriteString(rs.stdin, command)
	if err != nil {
		fmt.Println("Could not write to R process", err)
		return []byte{}, err
	}

	_, err = rs.getResult(key, filename)
	if err != nil {
		return []byte{}, errors.Wrap(err, "r: could not get results"+key)
	}

	info, err := os.Stat(filename)
	if err != nil {
		return []byte{}, errors.Wrap(err, "r: could not read fileinfo"+filename)
	}

	// wait until file is written completely
	size := int64(-1)
	for size != info.Size() {
		info, err = os.Stat(filename)
		if err != nil {
			return []byte{}, errors.Wrap(err, "could not read results file")
		}
		size = info.Size()
		time.Sleep(10 * time.Millisecond)
	}

	var b []byte
	b, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, errors.Wrap(err, "could not read results file")
	}

	if format == "json" {
		err = errors.New("")
		for err != nil {
			var js interface{}
			err = json.Unmarshal(b, &js)
			b, _ = ioutil.ReadFile(filename)
		}
	}

	return b, nil
}

func NewSession(id int) (*Session, error) {
	cmd := exec.Command("R", "-q", "--vanilla", "--slave")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	go func(cmd *exec.Cmd) {
		err = cmd.Run()
	}(cmd)

	// ensures that the R session has started before returning to caller
	// TODO: fix nasty hack.
	time.Sleep(200 * time.Millisecond)

	return &Session{cmd, stdin, stdout, stderr, id}, nil
}

func InitServer(numWorkers int, dir string) (Server, error) {
	rootDir = dir

	// clean tmp dir before starting up
	err := os.RemoveAll(rootDir)
	if err != nil {
		return Server{}, err
	}

	Sessions := make(chan *Session, numWorkers)

	for i := 0; i < numWorkers; i++ {
		rs, err := NewSession(i)
		if err != nil {
			fmt.Println("Could not start R session")
			return Server{}, err
		}
		Sessions <- rs

	}

	keys := make(chan string)

	return Server{keys, Sessions}, nil

}

func (rs *Session) InstalledPackages() ([]byte, error) {
	pkg := "utils"
	fun := "installed.packages"
	args := ""

	s, err := rs.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("Could not get installed packages", err)
		return nil, err
	}

	pkg = "base"
	fun = "as.data.frame"
	args = "x=" + s + ",row.names=make.names(installed.packages(), unique=TRUE)"

	s, err = rs.Call(pkg, fun, args)

	return rs.Get(s, "json")

}

type Call struct {
	Package   string
	Function  string
	Arguments string
}

func (c Call) cacheKey() string {
	return c.Package + "::" + c.Function + "(" + c.Arguments + ")"
}
