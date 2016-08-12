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
var rootDir = "/Users/bjorn/Dropbox/go/src/github.com/fjukstad/kvik/r/tmp/kvikr"

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

	io.WriteString(rs.stdin, "rm(list=ls());")
	io.WriteString(rs.stdin, "setwd(\""+wd+"\");")
	io.WriteString(rs.stdin, varName+"="+pkg+"::"+fun+"("+args+")"+";")
	io.WriteString(rs.stdin, "save.image();")
	io.WriteString(rs.stdin, varName+"\n")

	res, err := rs.getResult(key, wd+"/.RData")
	if err != nil {
		fmt.Println("Call failed")
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
		case <-time.After(1 * time.Second):
			return "", errors.New("R Call timeout. Check your syntax!")

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

	if format == "json" {
		_, err = io.WriteString(rs.stdin, "js=jsonlite::toJSON("+varName+"); write(js, file='output.json');\n")
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

	return &RSession{cmd, stdin, stdout, stderr, id}, nil
}

func InitServer(numWorkers int) (RServer, error) {
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

func main() {
	s, err := InitServer(15)
	if err != nil {
		return
	}
	for i := 0; i < 10; i++ {
		key, err := s.Call("stats", "rnorm", "n=100")
		if err != nil {
			fmt.Println("Call fail", err)
		} else {
			fmt.Println("Call success!", key)
		}
	}

	for i := 0; i < 50; i++ {
		key, err := s.Call("stats", "rnorm", "n=100")
		fmt.Println("keys", key, err)
		res, err := s.Get(key, "json")
		fmt.Println("results:", string(res), err)
	}
}
