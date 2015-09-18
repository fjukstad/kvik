package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

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
		// if it's not dependent on anything we
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

type Result struct {
	m     *sync.Mutex
	c     *sync.Cond
	Key   string
	Error error
}

func (p *Pipeline) WorkMagic() ([]*Result, error) {

	resultMap := make(map[string]*Result, 0)

	for _, stage := range p.Stages {
		m := sync.Mutex{}
		c := sync.NewCond(&m)
		resultMap[stage.Name] = &Result{&m, c, "", nil}
	}

	for _, stage := range p.Stages {
		deps := stage.GetDependencies()
		r := resultMap[stage.Name]

		if len(deps) == 0 {
			go func(r *Result, stage *Stage) {
				r.m.Lock()
				key, err := p.ExecuteStage(stage)
				r.Key = key
				r.Error = err
				//resultMap[stage.Name] = &Result{l, key, err}
				r.m.Unlock()
				r.c.Broadcast()

			}(r, stage)
		} else {
			m := make(chan bool, len(deps)-1)

			for _, dep := range deps {
				r := resultMap[dep]
				go func(r *Result, stage *Stage, dep string) {

					r.m.Lock()

					for r.Key == "" {
						r.c.Wait()
					}

					stage.ReplaceArg(dep, r.Key)
					m <- true
				}(r, stage, dep)
			}

			//		PARALLELIZE FOR LOOP:

			for i := 0; i < len(deps); i++ {
				<-m
			}

			go func(r *Result, stage *Stage) {
				r.m.Lock()

				r.Key, r.Error = p.ExecuteStage(stage)

				r.m.Unlock()
				r.c.Broadcast()

			}(r, stage)

		}
	}

	results := []*Result{}
	for _, stage := range p.Stages {
		r := resultMap[stage.Name]
		if r.Key == "" {
			r.m.Lock()
		}
		for r.Key == "" {
			r.c.Wait()
		}
		results = append(results, r)
	}

	return results, nil
}

func (s *Stage) ReplaceArg(oldarg string, newarg string) {
	for i, arg := range s.Arguments {
		if strings.Contains(arg, oldarg) {
			s.Arguments[i] = newarg
			return
		}
	}
}

func (p *Pipeline) ExecuteStage(stage *Stage) (string, error) {

	args := stage.GetArguments()
	fun := stage.GetFullFunctionName()

	s, err := p.Kompute.Call(fun, args)

	if err != nil {
		s, err = p.Kompute.Call(fun, args)

		if err != nil {
			fmt.Println("ERROR", err)
			return "", err
		}
	}
	return s.Key, nil
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

		args := stage.GetArguments()
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

func (s *Stage) Print() {
	fmt.Println("\tName:", s.Name)
	fmt.Println("\tPackage", s.Package)
	fmt.Println("\tFunction", s.Function)
	fmt.Println("\tArguments:")
	for k, v := range s.Arguments {
		fmt.Println("\t\t", k, "\t", v)
	}
	if s.Session != nil {
		fmt.Println("\tSession: ", s.Session.Key)
		fmt.Println("\tURL: /ocpu/tmp/" + s.Session.Key + "/R/")
	}

}
