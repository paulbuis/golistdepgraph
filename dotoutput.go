package main

import (
	"fmt"
	"io"
	"sort"
)

type IdState struct {
	nextId int
	ids    map[string]int
}

func (idState *IdState) getId(name string) int {
	id, ok := idState.ids[name]
	if !ok {
		id = idState.nextId
		idState.nextId++
		idState.ids[name] = id
	}
	return id
}

func dotOutput(pkgs map[string]JsonObject, dc DepContext, out io.Writer) {
	//importMap := make(map[string]int)
	//for index, pkg := range pkgs {
	//	importMap[pkg.GetString("ImportPath")] = index
	//}
	importPaths := []string{}
	for importPath, _ := range pkgs {
		importPaths = append(importPaths, importPath)
	}
	sort.Strings(importPaths)
	idState := &IdState{ids:make(map[string]int)}

	fmt.Fprintln(out, "digraph G {")
	for _, importPath := range importPaths {
		pkg := pkgs[importPath]
		pkgId := idState.getId(importPath)

		if dc.IsIgnored(pkg) {
			continue
		}

		var color string
		if pkg.GetBool("Goroot") {
			color = "palegreen"
		} else if len(pkg.GetStringSlice("CgoFiles")) > 0 {
			color = "darkgoldenrod1"
		} else {
			color = "paleturquoise"
		}
		var fontColor string
		if pkg.GetBool("Incomplete") {
			fontColor = "red"
		 } else if pkg.GetBool("Stale") {
			fontColor = "blue"
		 } else {
			 fontColor = "black"
		 }
		// what about pkg.GetBool("Incomplete"), pkg.GetBool("Stale"), pkg.GetString("StaleReason")?

		fmt.Fprintf(out, "_%d [label=\"%s\" style=\"filled\" color=\"%s\" fontcolor=\"%s\"];\n", pkgId, pkg.GetString("ImportPath"), color, fontColor)

		if pkg.GetBool("Goroot") && !dc.DelveGoroot {
			continue
		}

		for _, imp := range pkg.GetStringSlice("Imports") {
			impPkg, ok := pkgs[imp]
			if !ok || dc.IsIgnored(impPkg) {
				continue
			}
			impId := idState.getId(imp)
			fmt.Fprintf(out, "_%d -> _%d;\n", pkgId, impId)
		}
	}
	fmt.Fprintln(out, "}")
}
