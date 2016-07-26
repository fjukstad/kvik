package r

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

// Installs the given packages on the server. packages is a space
// separated list with the packages.
func installPackages(packages string) error {
	cmd := exec.Command("R", "-q", "-e", "install.packages(c("+packages+"),  repos='http://cran.us.r-project.org'), dependencies=TRUE")

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

func newCall() (string, string, error) {
	h := md5.New()
	k := h.Sum([]byte(strconv.Itoa(r.Int())))

	hash := hex.EncodeToString(k)
	key := ".s" + hash

	wd := rootDir + "/" + key
	err := os.MkdirAll(wd, 0755)
	if err != nil {
		return "", "", err
	}

	err = os.Chdir(wd)

	if err != nil {
		return "", "", err
	}

	return key, wd, nil

}

// Executes function. Returns tmp key for use in Get
func Call(pkg, fun, args string) (*Session, error) {

	key, wd, err := newCall()
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

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	var output string
	if err != nil {
		output = stderr.String()
		// remove kvik tmp key
		output = strings.Replace(output, varName+"=", "", -1)
	} else {
		output = stdout.String()
	}

	return &Session{key, wd, output, command}, err
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

func ScriptHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("script handler")
	_, wd, err := newCall()
	if err != nil {
		fmt.Println("new call fail:", err)
		w.Write([]byte(err.Error()))
		return
	}

	file, _, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Println("workdir", wd)

	fileLoc := wd + "/script.R"

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

	output, err := executeScript(wd)
	if err != nil {
		w.Write([]byte("ERROR, could not execute script. Output:\n" + output))
		return
	}

	cmd := exec.Command("zip", "-r", "-q", "output.zip", ".")
	cmd.Dir = wd

	var outbuffer bytes.Buffer
	cmd.Stdout = &outbuffer

	err = cmd.Run()

	if err != nil {
		fmt.Println("ERROR", outbuffer.String())
		w.Write([]byte("Could not zip output file"))
		return
	}

	http.ServeFile(w, r, "output.zip")

	return

}

func executeScript(wd string) (string, error) {

	filename := wd + "/script.R"

	cmd := exec.Command("R", "--save", "-f", filename)
	cmd.Dir = wd

	fmt.Println(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {

		errMessage := out.String()
		fmt.Println("ERROR", errMessage)
		return errMessage, err
	}

	return out.String(), nil
}

func InstalledPackages() ([]byte, error) {

	pkg := "utils"
	fun := "installed.packages"
	args := ""

	s, err := Call(pkg, fun, args)
	if err != nil {
		fmt.Println("Could not get installed packages")
		fmt.Println(s.Output)
		return nil, err
	}

	pkg = "base"
	fun = "as.data.frame"
	args = "x=" + s.Key + ",row.names=make.names(installed.packages(), unique=TRUE)"

	s, err = Call(pkg, fun, args)
	if err != nil {
		fmt.Println("Could not get installed packages")
		fmt.Println(s.Output)
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
