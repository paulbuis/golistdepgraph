package main

import (
	"fmt"
)

// analogous to godepgraph's processPackage

func WalkDepGraph(dir string, pkgImportPath string, dc DepContext, pkgMap map[string]JsonObject) error {
	if dc.Ignored[pkgImportPath] {
		return nil
	}

	pkg, err := JsonImmediateDep(dir, pkgImportPath)
	if err != nil {
		return fmt.Errorf("failed to import %s, %s", pkgImportPath, err)
	}

	if dc.IsIgnored(pkg) {
		return nil
	}

	pkgMap[pkg.GetString("ImportPath")] = pkg

	if (pkg.GetBool("Goroot")) && !dc.DelveGoroot {
		return nil
	}

	for _, imp := range getImports(dc, pkg) {
		if _, ok := pkgMap[imp]; !ok {
			if err := WalkDepGraph(dir, imp, dc, pkgMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func getImports(dc DepContext, pkg JsonObject) []string {
	allImports := pkg.GetStringSlice("Imports")
	if dc.IncludeTests {
		allImports = append(allImports, pkg.GetStringSlice("TestImports")...)
		allImports = append(allImports, pkg.GetStringSlice("XTestImports")...)
	}
	var imports []string
	found := make(map[string]struct{})
	pkgImportPath := pkg.GetString("ImportPath")
	for _, imp := range allImports {
		if imp == pkgImportPath {
			// avoiding self-reference
			continue
		}
		if _, ok := found[imp]; ok {
			// skipping repeated packges contained in allImports
			continue
		}
		found[imp] = struct{}{}
		imports = append(imports, imp)
	}
	return imports
}
