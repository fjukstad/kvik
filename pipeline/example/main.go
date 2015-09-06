package main

import (
	"fmt"

	"github.com/fjukstad/kvik/kompute"
	"github.com/fjukstad/kvik/pipeline"
)

func main() {
	addr := "192.168.99.100:8004"
	username := "user"
	password := "password"

	k := kompute.NewKompute(addr, username, password)

	p := pipeline.NewPipeline("stress-pipe", k)

	// --------------- LOAD DATA ---------------- //
	name := "loaddata"
	function := "LoadData"
	pkg := "stress"
	argnames := []string{"filename", "bad.probes.filename"}
	args := []string{
		"/home/stallo/src/stress/datasets/StressTest_rev8.RData",
		"/home/stallo/src/stress/datasets/hla_hist.txt"}

	s := pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	// --------------- LOAD BACKGROUND DATA ---------------- //
	name = "loadbackgrounddata"
	function = "LoadBackgroundData"
	pkg = "stress"
	argnames = []string{"dataset", "filename"}
	args = []string{
		"from:loaddata",
		"/home/stallo/src/stress/datasets/cc-status-stressv1.csv"}

	s = pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	// --------------- FILTER PROBES ---------------- //
	name = "filterprobes"
	function = "FilterProbes"
	pkg = "stress"
	argnames = []string{"dataset", "pVal", "fVal"}
	args = []string{
		"from:loaddata",
		"0.01",
		"0.25"}

	s = pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	// --------------- ANNOTATE PROBES ---------------- //
	name = "annotateprobes"
	function = "AnnotateProbes"
	pkg = "stress"
	argnames = []string{"dataset"}
	args = []string{
		"from:filterprobes"}

	s = pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	// --------------- Group 0 t-test ---------------- //
	name = "group0t-test"
	function = "CaseControlTTest"
	pkg = "stress"
	argnames = []string{"dataset", "clinical", "c_status"}
	args = []string{
		"from:annotateprobes",
		"from:loadbackgrounddata",
		"0",
	}

	s = pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	// --------------- Group 1 t-test ---------------- //
	name = "group1t-test"
	function = "CaseControlTTest"
	pkg = "stress"
	argnames = []string{"dataset", "clinical", "c_status"}
	args = []string{
		"from:annotateprobes",
		"from:loadbackgrounddata",
		"1",
	}

	s = pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	// --------------- CC Fold Change t-test ---------------- //
	name = "ccfct-test"
	function = "CaseControlFoldChangeTTest"
	pkg = "stress"
	argnames = []string{"dataset", "clinical"}
	args = []string{
		"from:annotateprobes",
		"from:loadbackgrounddata",
	}

	s = pipeline.NewStage(name, function, pkg, argnames, args)
	p.AddStage(s)

	p.Run()

	fmt.Println("done...")

	p.Save()

}
