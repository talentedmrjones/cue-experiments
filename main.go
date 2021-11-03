package main

import (
	"fmt"
	"os"
	"regexp"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
)

func main() {

	buildInstances := getBuildInstances(nil)
	plan := getValue(buildInstances)

	context := plan.LookupPath(cue.ParsePath("context"))

	context.Walk(func(contextField cue.Value) bool {

		fmt.Printf("%s -> %s\n", contextField.Path(), contextField.IncompleteKind())

		sources := contextField.Split()
		for _, src := range sources {
			fmt.Printf("%s -> %+v\n", src.Source().Pos().Filename(), src.IsConcrete())
		}
		return true
	}, nil)

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
		args = append(args, "./plans/dev")
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

func setByPackage(v cue.Value) bool {
	fmt.Printf("%s -> %s\n", v.Path(), v.Kind())
	if v.IsConcrete() && fmt.Sprint(v.Kind()) == "string" {
		// fmt.Printf("%s\n", v.Value())
		filename := v.Source().Pos().File().Name()
		match, _ := regexp.MatchString("cue.mod/(pkg|usr)", filename)
		if match {
			fmt.Printf("CANNOT SET CONCRETE VALUE IN PACKAGE %s\n", filename)
		} else {
			fmt.Printf("%s\n", v.Value())
		}
	} else if v.Kind() == cue.BottomKind {
		_, defaultExisted := v.Default()
		filename := v.Source().Pos().File().Name()
		match, _ := regexp.MatchString("cue.mod/(pkg|usr)", filename)
		if defaultExisted && match {
			fmt.Printf("CANNOT SET DEFAULT VALUE IN PACKAGE %s\n", filename)
		}
	}
	return true
}
