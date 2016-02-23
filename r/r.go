package R

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	if err != nil {
		return nil, err
	}

	return &Session{key, wd, out.String(), command}, nil
}

func Get(key, format string) ([]byte, error) {

	dir := rootDir + "/" + key
	err := os.Chdir(dir)

	varName := strings.TrimPrefix(key, ".")

	if err != nil {
		return nil, err
	}

	extension := "." + format
	_, err = os.Stat("output" + extension)
	if err == nil {
		return ioutil.ReadFile(dir + "/output" + extension)
	}

	var command string
	if format == "csv" {
		command = "write.csv(" + varName + ", sep=',', file='output" + extension + "')"
	} else if format == "json" {
		command = "js=jsonlite::toJSON(" + varName + "); write(js, file='output" + extension + "')"
	} else if format == "pdf" {
		return ioutil.ReadFile(dir + "/Rplots.pdf")
	} else {
		return nil, errors.New("Unknown format")
	}

	cmd := exec.Command("R", "--save", "-q", "-e", command)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(dir + "/output" + extension)

}
