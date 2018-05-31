package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func JsonImmediateDep(path string, pkgName string) (JsonObject, error) {
	pkgSrc := pkgName
	cmd := exec.Command("go", "list", "-e", "-json", pkgSrc)
	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return NewJsonObject(), err
	}
	errStr := string(stderr.Bytes())
	if len(errStr) > 0 {
		fmt.Printf("stderr contained: %s\n", errStr)
	}
	js := NewJsonSeq(stdout.Bytes())
	return js[0], nil
}