package internal

import (
	"bytes"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildPlugin(output string, dir string) error {
	b, e := build.ImportDir(dir, build.FindOnly)
	if e != nil {
		return e
	}

	info, e := os.Stat(dir)
	if e != nil {
		return e
	}

	if !info.IsDir() {
		return fmt.Errorf("%s not a directory", dir)
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", fmt.Sprintf("%s/%s.so", output, info.Name()), ".")
	cmd.Dir = b.Dir

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	e = cmd.Run()

	if len(stderr.String()) > 0 {
		fmt.Printf("errors when building %s:\n%s\n", dir, stderr.String())
	}

	if len(stdout.String()) > 0 {
		fmt.Printf("when building %s:\n%s\n", dir, stdout.String())
	}

	if e != nil {
		return e
	}

	return nil
}

func BuildPlugins(output string, dirs []string) error {

	info, e := os.Stat(output)
	if e != nil {
		return e
	}

	if !info.IsDir() {
		return fmt.Errorf("output is not a directory")
	}

	o, e := filepath.Abs(output)
	if e != nil {
		return e
	}

	for _, dir := range dirs {
		e := BuildPlugin(o, dir)
		if e != nil {
			return e
		}
	}

	return nil
}
