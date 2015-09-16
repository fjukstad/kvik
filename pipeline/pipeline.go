package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fjukstad/kvik/kompute"
)

type Stage struct {
	Name      string
	Package   string
	Function  string
	Arguments map[string]string
	Depends   []Stage
	Session   *kompute.Session
}

type Pipeline struct {
	Name    string
	Kompute *kompute.Kompute
	Stages  []*Stage
}

func NewPipeline(name string, k *kompute.Kompute) Pipeline {
	p := Pipeline{name, k, nil}
	return p
}

func (p *Pipeline) AddStage(s Stage) {
	p.Stages = append(p.Stages, &s)
}

func NewStage(name, function, pkg string, argnames, args []string) Stage {

	argmap := make(map[string]string, 0)

	for i, argname := range argnames {
		arg := args[i]
		argmap[argname] = arg
	}

	s := Stage{name, pkg, function, argmap, []Stage{}, nil}

	return s
}

func (p *Pipeline) RunParallel() ([]*Pipeline, error) {

	pipeMap := map[string]*Pipeline{}
	pipes := []*Pipeline{}

	i := 0

	for _, stage := range p.Stages {

		deps := stage.GetDependencies()
		// if it's not dependent on anything we can place it
		// in a new pipeline.
		if len(deps) == 0 {
			pipe := NewPipeline(p.Name+"-par-"+strconv.Itoa(i), p.Kompute)
			pipe.AddStage(*stage)
			pipeMap[stage.Name] = &pipe
			i++

			pipes = append(pipes, &pipe)

		} else {
			pipe := pipeMap[deps[0]]
			fmt.Println(pipe)
			pipe.AddStage(*stage)
			pipeMap[stage.Name] = pipe
			fmt.Println("Pipeline", len(pipe.Stages))
		}
	}

	for _, pipe := range pipes {
		fmt.Println(pipe.Name, len(pipe.Stages))
		pipe.Run()
	}

	return pipes, nil
}

func (s Stage) GetDependencies() []string {
	deps := []string{}
	for _, arg := range s.Arguments {
		if strings.Contains(arg, "from:") {
			argname := strings.Split(arg, "from:")[1]
			deps = append(deps, argname)
		}
	}

	return deps
}

func (p *Pipeline) GetResult(format string) string {
	lastStage := p.Stages[len(p.Stages)-1]
	res, _ := lastStage.Session.GetResult(p.Kompute, format)
	return res
}

func (p *Pipeline) Run() error {

	for _, stage := range p.Stages {

		args := stage.GetArguments(p)
		fun := stage.GetFullFunctionName()

		s, err := p.Kompute.Call(fun, args)
		if err != nil {
			s, err = p.Kompute.Call(fun, args)
			if err != nil {
				fmt.Println("ERROR", err)
				return err
			}
		}

		stage.Session = s
	}

	return nil
}

func (p *Pipeline) Save() error {
	file, err := os.OpenFile(p.Name+".json", os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *Stage) GetFullFunctionName() string {
	return s.Package + "/R/" + s.Function
}

func (s *Stage) GetArguments(p *Pipeline) string {

	str := ""
	i := 0

	for k, v := range s.Arguments {

		// fix dependency. get session key from "from:" stage
		if strings.Contains(v, "from:") {
			stageName := strings.Split(v, "from:")[1]
			for _, stage := range p.Stages {
				if stage.Name == stageName {
					v = stage.Session.Key
				}
			}
		}

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

func (s *Stage) Print() {
	fmt.Println("\tName:", s.Name)
	fmt.Println("\tPackage", s.Package)
	fmt.Println("\tFunction", s.Function)
	fmt.Println("\tArguments:")
	for k, v := range s.Arguments {
		fmt.Println("\t\t", k, "\t", v)
	}
	fmt.Println("\tSession: ", s.Session.Key)
	fmt.Println("\tURL: /ocpu/tmp/" + s.Session.Key + "/R/")

}
