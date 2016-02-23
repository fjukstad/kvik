// Package for building an execution pipeline and executing it.
package pipeline

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/fjukstad/kvik/gopencpu"
)

// A pipeline stage.
// Name: A unique name for each stage.
// Function: Function name, e.g. plot.
// Package: Which statistical package the function comes from, e.g. graphics for plot().
// Arguments: A name:value map for arguments to the function.
// Depends: List of stages it depends on.
// Session: Reference to the OpenCPU session used to compute the results from this stage.
// Output: String for storing output from pipeline stage execution.
type Stage struct {
	Name      string            "name,omitempty"
	Package   string            "package,omitempty"
	Function  string            "function,omitempty"
	Arguments map[string]string "arguments,omitempty"
	Depends   []string          "depends,omitempty"
	Session   *gopencpu.Session "session,omitempty"
	Output    string            "output,omitempty"
}

// The pipeline. A name, reference to an opencpu server and list of pipeline
// stages.
type Pipeline struct {
	Name      string              "name,omitempty"
	GoOpenCPU *gopencpu.GoOpenCPU "gopencpu,omitempty"
	Stages    []*Stage            "stages,omitempty"
}

// Set up new pipeline with the given name and reference to gopencpu server
func NewPipeline(name string, k *gopencpu.GoOpenCPU) Pipeline {
	p := Pipeline{name, k, nil}
	return p
}

// Appends a pipeline stage to to the end of the pipeline.
func (p *Pipeline) AddStage(s Stage) {
	p.Stages = append(p.Stages, &s)
}

// Creates a new pipeline stage. Argnames is an array of argument names, args
// are the values for each of the arguments.  Assumes that argnames and args are
// of the same length.
func NewStage(name, function, pkg string, argnames, args []string) Stage {

	argmap := make(map[string]string, 0)

	for i, argname := range argnames {
		arg := args[i]
		argmap[argname] = arg
	}

	s := Stage{name, pkg, function, argmap, []string{}, nil, ""}

	return s
}

// Imports a pipeline from a pipeline description in yaml. It should follow this
// format:
//
//	 name: PIPELINE NAME
//	 stages:
//	 - name: STAGE NAME
//	   package: PACKAGE NAME
//	   function: FUNCTION NAME
//	   arguments:
//	     ARGUMENT NAME: ARGUMENT VALUE
//	     ...

func ImportPipeline(filename string) (Pipeline, error) {
	p := Pipeline{}

	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return Pipeline{}, err
	}

	err = yaml.Unmarshal(in, &p)
	if err != nil {
		return Pipeline{}, err
	}
	return p, nil

}

// For storing the result of a pipeline stage. Key can be used to retrieve
// plots/results later using the gopencpu package.
type Result struct {
	m     *sync.Mutex
	c     *sync.Cond
	Key   string
	Error error
}

// Exectues a pipeline. Uses a go routine per pipeline stage making it possible
// to execute multiple stages simultaneously. Returns a list of results, one per
// stage.
func (p *Pipeline) Run() ([]*Result, error) {

	resultMap := make(map[string]*Result, 0)
	done := make(chan bool, len(p.Stages))

	for _, stage := range p.Stages {
		m := sync.Mutex{}
		c := sync.NewCond(&m)
		resultMap[stage.Name] = &Result{&m, c, "", nil}

		deps := stage.GetDependencies()
		r := resultMap[stage.Name]
		stage.Depends = deps

		go func(r *Result, stage *Stage, deps []string) {
			if len(deps) != 0 {
				for _, dep := range deps {
					r := resultMap[dep]
					r.m.Lock()

					for r.Key == "" {
						r.c.Wait()
					}

					stage.ReplaceArg(dep, r.Key)
					r.m.Unlock()
				}
			}
			r.m.Lock()
			for r.Key == "" {
				r.Key, r.Error = p.ExecuteStage(stage)
			}
			r.c.Broadcast()
			r.m.Unlock()

			done <- true
		}(r, stage, deps)
	}

	for i := 0; i < len(p.Stages); i++ {
		<-done
	}

	results := []*Result{}
	for _, stage := range p.Stages {
		r := resultMap[stage.Name]
		results = append(results, r)
	}

	return results, nil
}

func (p *Pipeline) RunSequential() ([]*Result, error) {
	resultMap := make(map[string]*Result, 0)

	for _, stage := range p.Stages {
		resultMap[stage.Name] = &Result{nil, nil, "", nil}

		deps := stage.GetDependencies()
		r := resultMap[stage.Name]

		for _, dep := range deps {
			r := resultMap[dep]
			stage.ReplaceArg(dep, r.Key)
		}

		for r.Key == "" {
			r.Key, r.Error = p.ExecuteStage(stage)
		}

	}

	results := []*Result{}
	for _, stage := range p.Stages {
		r := resultMap[stage.Name]
		results = append(results, r)
	}

	return results, nil

}

// Replace any "from:stage-name" argument values into opencpu references
func (s *Stage) ReplaceArg(oldarg string, newarg string) {
	for i, arg := range s.Arguments {
		if strings.Contains(arg, oldarg) {
			s.Arguments[i] = strings.Replace(s.Arguments[i], "from:"+oldarg, newarg, -1)
			return
		}
	}
}

// Executes a pipeline stage.
func (p *Pipeline) ExecuteStage(stage *Stage) (string, error) {

	args := stage.GetArguments()
	fun := stage.GetFullFunctionName()

	s, err := p.GoOpenCPU.Call(fun, args)

	if err != nil {
		s, err = p.GoOpenCPU.Call(fun, args)

		if err != nil {
			fmt.Println("ERROR", err)
			return "", err
		}
	}

	fmt.Println("Stage", stage.Name, "completed")

	stage.Session = s
	return s.Key, nil
}

// Get the final result from the last stage in a pipeline.
func (p *Pipeline) GetResult(format string) string {
	lastStage := p.Stages[len(p.Stages)-1]
	res, _ := lastStage.Session.GetResult(p.GoOpenCPU, format)
	return res
}

// Stores a pipeline as a pipeline description yaml file.
func (p *Pipeline) Save() error {
	file, err := os.OpenFile(p.Name+".yaml", os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// Returns the full function name to be used with the gopencpu package.
func (s *Stage) GetFullFunctionName() string {
	return s.Package + "/R/" + s.Function
}

// Create a string of arguments for usage by the Call method in the gopencpu
// package
func (s *Stage) GetArguments() string {
	str := ""
	i := 0
	for k, v := range s.Arguments {
		if strings.Contains(v, "/") || strings.Contains(v, ".") {
			v = "\"" + v + "\""
		}
		str = str + k + "=" + v
		if i < len(s.Arguments)-1 {
			str = str + "&"
		}
		i++
	}

	return str
}

// Returns all dependencies for a stage.
func (s Stage) GetDependencies() []string {
	deps := []string{}
	for _, arg := range s.Arguments {
		if strings.Contains(arg, "from:") {
			args := strings.Split(arg, "from:")
			var argname string

			// if argument is list of from: s
			if len(args) > 2 {
				for _, a := range args {
					if len(a) > 1 {
						a = strings.TrimRight(a, ",")
						a = strings.TrimRight(a, "]")
						deps = append(deps, a)
					}
				}
			} else {
				argname = strings.Split(arg, "from:")[1]
				argname = strings.TrimRight(argname, "]")
				deps = append(deps, argname)
			}
		}
	}

	return deps
}

// Print a pipeline stage.
func (p *Pipeline) Print() {
	for _, stage := range p.Stages {
		if stage.Session != nil {
			res, _ := p.GoOpenCPU.Get(stage.Session.Key, "")
			stage.Output = string(res)
		}
		stage.Print()
		fmt.Println()
	}
}

// Return pipeline results, i.e. result of last pipeline stage.
func (p *Pipeline) Results(resultType string) (string, error) {
	stage := p.Stages[len(p.Stages)-1]

	if resultType == "png" || resultType == "pdf" {

		a := "/Users/bjorn/Dropbox/go/src/github.com/fjukstad/kvik/pipeline/example/x-y/"

		return "", stage.Session.DownloadPlot(resultType, a+p.Name+"."+resultType)
	}

	return stage.Session.GetResult(p.GoOpenCPU, resultType)
}

// Print a pipeline stage.
func (s *Stage) Print() {
	fmt.Println("\tName:", s.Name)
	fmt.Println("\tPackage", s.Package)
	fmt.Println("\tFunction", s.Function)
	fmt.Println("\tDepends on:", s.Depends)
	fmt.Println("\tArguments:")
	for k, v := range s.Arguments {
		fmt.Println("\t\t", k, "\t", v)
	}
	if s.Session != nil {
		fmt.Println("\tSession: ", s.Session.Key)
		fmt.Println("\tURL: /ocpu/tmp/" + s.Session.Key + "/R/")
		fmt.Println("\tOutput:\n", s.Output)
	}
}
