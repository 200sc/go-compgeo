package main

// Build cross-compiles the demo package on a small set of
// OS and architecture pairs. It should be generalized to
// take in variable output and package names, and varible
// sets of os-arch pairs, then split off into a separate
// project.

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

var (
	osArchPairs = [][2]string{
		{"linux", "386"},
		{"linux", "amd64"},
		{"linux", "arm"},
		{"linux", "arm64"},
		{"windows", "386"},
		{"windows", "amd64"},
	}
	packageName = "github.com/200sc/go-compgeo/demo"
	outputName  = "pl-demo"
	verbose     = true
)

func main() {
	goos := os.Getenv("GOOS")
	goarch := os.Getenv("GOARCH")
	for _, pair := range osArchPairs {
		os.Setenv("GOOS", pair[0])
		os.Setenv("GOARCH", pair[1])
		var out bytes.Buffer
		if verbose {
			fmt.Println("Running: go build -o", outputName+"_"+pair[0]+pair[1], packageName)
		}
		cmd := exec.Command("go", "build", "-o", outputName+"_"+pair[0]+pair[1], packageName)
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		if verbose && out.Len() != 0 {
			fmt.Printf("%s\n", out.String())
		}
	}
	os.Setenv("GOOS", goos)
	os.Setenv("GOARCH", goarch)
}
