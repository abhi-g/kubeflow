// Package main provides an application to run load tests
// against the Kubeflow deployer application.
package main

import (
	"flag"
	"fmt"

	"github.com/kubeflow/kubeflow/bootstrap/cmd/load-test/projects"
)

var flags = struct {
	numProjects int
	runTest     bool
}{}

func initFlags() {
	flag.IntVar(&flags.numProjects, "create-projects", 0,
		"number of projects to create for deployment load test")

	flag.BoolVar(&flags.runTest, "run", false,
		"run the load test when set, otherwise skip")
}

func main() {
	fmt.Println("Hello load test")

	initFlags()
	flag.Parse()

	fmt.Println("num projects:", flags.numProjects)
	fmt.Println("run test:", flags.runTest)

	if flags.numProjects > 0 {
		projects.CreateAllProjects(flags.numProjects)
	}

	if flags.runTest {
		// run the load test
	}
}
