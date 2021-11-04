package main

import (
	"fmt"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
)

// this will take the path at arg[1] and look for cue.mod in parent directories recursively

func main() {

	var cwd, projectRoot string

	if len(os.Args) > 1 {
		cwd = os.Args[1]
	} else {
		cwd = "./"
	}

	path, _ := filepath.Abs(cwd)

	// traverse the directory tree starting from PWD going up to successive parents
	for {
		// look for the cue.mod filder
		if _, err := os.Stat(path + "/cue.mod"); !os.IsNotExist(err) {
			projectRoot = path
			break // found it!
		}
		path, _ = filepath.Abs(filepath.Dir(path))
		if path == string(os.PathSeparator) {
			break
		}
	}

	if projectRoot == "" {
		fmt.Println("Cannot determine project root. No cue.mod found.")
	}

	fmt.Printf("%s", projectRoot)

}

func getBuildInstances(args []string) []*build.Instance {
	// const syntaxVersion = -1000 + 13

	config := load.Config{
		Package: "main",
		Context: build.NewContext(
			build.ParseFile(func(name string, src interface{}) (*ast.File, error) {
				return parser.ParseFile(name, src,
					parser.FromVersion(parser.Latest),
					parser.ParseComments,
				)
			})),
	}

	if len(args) < 1 {
		args = append(args, "./")
	}

	buildInstances := load.Instances(args, &config)

	return buildInstances
}

func getValue(buildInstances []*build.Instance) cue.Value {
	ctx := cuecontext.New()

	values, buildInstancesErr := ctx.BuildInstances(buildInstances)
	if buildInstancesErr != nil {
		fmt.Printf("%s", buildInstancesErr)
		os.Exit(1)
	}
	return values[0]
}
