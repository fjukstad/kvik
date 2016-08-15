package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var rootDir string
var timeout int64 = 4

type RSession struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	id     int
}

type RServer struct {
	keys     chan string
	sessions chan *RSession
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s RServer) Call(pkg, fun, args string) (string, error) {
	session := <-s.sessions
	res, err := session.Call(pkg, fun, args)
	s.sessions <- session
	return res, err
}

func (s RServer) Get(key, format string) ([]byte, error) {
	session := <-s.sessions
	res, err := session.Get(key, format)
	s.sessions <- session
	return res, err
}

func (rs *RSession) Call(pkg, fun, args string) (string, error) {

	key := "." + randSeq(5)
	wd := rootDir + "/" + key
	varName := strings.TrimPrefix(key, ".")

	err := os.MkdirAll(wd, 0755)
	if err != nil {
		return "", err
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
		fmt.Println("io write string fail")
	}

	fmt.Println("ran command", command)

	res, err := rs.getResult(key, wd+"/.RData")
	if err != nil {
		_, err = io.WriteString(rs.stdin, command)
		if err != nil {
			fmt.Println("io write fail")
			return "", err
		}
		res, err = rs.getResult(key, wd+"/.RData")
	}
	return res, err
}

func (rs *RSession) getResult(key, filename string) (string, error) {
	keys := make(chan string)
	errCh := make(chan string)

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

	for {
		select {
		// call successful no errors
		case k := <-keys:
			return k, nil
		// Execution of R code went wrong. Restart a new session
		case o := <-errCh:
			rs.cmd.Process.Kill()
			newSession, _ := NewSession(rs.id)
			rs.cmd = newSession.cmd
			rs.stdin = newSession.stdin
			rs.stdout = newSession.stdout
			rs.stderr = newSession.stderr
			rs.id = newSession.id
			return "", errors.New(o)
		case <-time.After(time.Duration(timeout) * time.Second):
			return "", errors.New("R Call timeout.")

		}
	}

	return key, nil
}

func (rs *RSession) Get(key, format string) ([]byte, error) {
	wd := rootDir + "/" + key

	cmd := "rm(list=ls());" + "setwd(\"" + wd + "\");" + "load(\".RData\");\n"

	_, err := io.WriteString(rs.stdin, cmd)
	if err != nil {
		fmt.Println("Could not set up for get")
		return []byte{}, nil
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
			fmt.Println("Could not read pdf file", err)
			return nil, err
		}
		// if file contains magic end return it, if not wait
		// for R to complete writing the file.
		for !strings.Contains(string(file), "%%EOF") {
			time.Sleep(500 * time.Millisecond)
			file, err = ioutil.ReadFile(wd + "/Rplots.pdf")
			if err != nil {
				fmt.Println("Could not read pdf file", err)
				return nil, err
			}
		}
		return file, err
	} else if format == "png" {
		//cmd := exec.Command("pdftoppm", "-png", wd+"/Rplots.pdf", "plot.png")
		cmd := exec.Command("convert", wd+"/Rplots.pdf", wd+"/plot-1.png")

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Could not convert rplot to png:", stderr.String())
			fmt.Println("Check that you have Xpdf installed.")
			return nil, err
		}

		return ioutil.ReadFile(wd + "/plot-1.png")
	} else {
		return nil, errors.New("Unknown format:" + format)
	}

	_, err = io.WriteString(rs.stdin, command)
	if err != nil {
		fmt.Println("Could not write to R process", err)
		return []byte{}, err
	}

	_, err = rs.getResult(key, filename)
	if err != nil {
		fmt.Println("Get failed")
		return []byte{}, err
	}
	var b []byte
	b, err = ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read results from file")
		return []byte{}, err
	}

	return b, nil
}

func NewSession(id int) (*RSession, error) {
	cmd := exec.Command("R", "-q", "--save")
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

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// ensures that the R session has started before returning to caller
	// TODO: fix nasty hack.
	time.Sleep(200 * time.Millisecond)

	return &RSession{cmd, stdin, stdout, stderr, id}, nil
}

func InitServer(numWorkers int, dir string) (RServer, error) {
	rootDir = dir
	RSessions := make(chan *RSession, numWorkers)

	for i := 0; i < numWorkers; i++ {
		rs, err := NewSession(i)
		if err != nil {
			fmt.Println("Could not start R session")
			return RServer{}, err
		}
		RSessions <- rs

	}

	keys := make(chan string)

	return RServer{keys, RSessions}, nil

}

func (rs *RSession) InstalledPackages() ([]byte, error) {
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

func main() {
	s, err := InitServer(10, "/Users/bjorn/Dropbox/go/src/github.com/fjukstad/kvik/r/tmp/kvikr")

	if err != nil {
		return
	}
	for i := 0; i < 10; i++ {
		key, err := s.Call("stats", "rnorm", "n=100")
		if err != nil {
			fmt.Println("Call fail", err)
			return
		} else {
			fmt.Println("Call success!", key)
		}
	}

	for i := 0; i < 50; i++ {
		key, err := s.Call("stats", "rnorm", "n=100")

		fmt.Println("keys", key, err)
		if err == nil {
			res, err := s.Call("graphics", "plot", "x="+key)
			fmt.Println("results:", string(res), err)
			if err != nil {
				fmt.Println(err)
				return
			}

			file, err := s.Get(res, "pdf")
			err = ioutil.WriteFile("pdf.pdf", file, 0777)

			fmt.Println("SHITPDF", err)
			fmt.Println(len(file))
		} else {
			fmt.Println("##############################################")
			return
		}
	}

	/*
		rs, err := NewSession(0)
		if err != nil {
			fmt.Println("ns err", err)
			return
		}
		pkg := "utils"
		fun := "installed.packages"
		args := ""

		sa, err := rs.Call(pkg, fun, args)
		if err != nil {
			fmt.Println("call failed", err)
		} else {
			res, err := rs.Get(sa, "json")
			fmt.Println(string(res))
			fmt.Println(err)
		}
	*/
}
