package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

type DepContext struct {
	Ignored         map[string]bool
	IgnoredPrefixes []string
	IgnoreStdlib    bool
	DelveGoroot     bool
	TagList         string
	IncludeTests    bool
	BuildTags       []string
}

func main() {
	var (
		ignoreStdlib   = flag.Bool("s", false, "ignore packages in the Go standard library")
		delveGoroot    = flag.Bool("d", false, "show dependencies of packages in the Go standard library")
		ignorePrefixes = flag.String("p", "", "a comma-separated list of prefixes to ignore")
		ignorePackages = flag.String("i", "", "a comma-separated list of packages to ignore")
		tagList        = flag.String("tags", "", "a comma-separated list of build tags to consider satisified during the build")
		includeTests   = flag.Bool("t", false, "include test packages")
	)

	flag.Parse()
	args := flag.Args()
	depContext := DepContext{
		IgnoreStdlib: *ignoreStdlib,
		DelveGoroot:  *delveGoroot,
		IncludeTests: *includeTests,
		Ignored: make(map[string]bool)}

	if *ignorePrefixes != "" {
		depContext.IgnoredPrefixes = strings.Split(*ignorePrefixes, ",")
	}
	if *ignorePackages != "" {
		for _, p := range strings.Split(*ignorePackages, ",") {
			depContext.Ignored[p] = true
		}
	}
	if *tagList != "" {
		depContext.BuildTags = strings.Split(*tagList, ",")
	}

	if len(args) != 1 {
		log.Fatal("need one package name to process")
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd: %s", err)
	}
	path := cwd
	pkgName := args[0]

	pkgMap := make(map[string]JsonObject)
	err = WalkDepGraph(path, pkgName, depContext, pkgMap)
	if err != nil {
		log.Fatal(err)
	}

	dotOutput(pkgMap, depContext, os.Stdout)
}

func hasPrefixes(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func (d DepContext) IsIgnored(pkg JsonObject) bool {
	importPath := pkg.GetString("ImportPath")
	return d.Ignored[importPath] || (pkg.GetBool("Goroot") && d.IgnoreStdlib) || hasPrefixes(importPath, d.IgnoredPrefixes)
}
