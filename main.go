package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/coreos/go-semver/semver"
)

func main() {
	getNext := flag.Bool("next", false, "Just print the next version")
	asNumber := flag.Bool("numeric", false, "Numeric form")
	flag.Parse()

	branch := getCurrentBranch()
	if !*getNext && (branch != "main" && branch != "master") {
		panic(fmt.Errorf(`error: must be in "master/main" branch, current branch: %q`, branch))
	}

	lastTag := getLastTag()

	version := strings.TrimPrefix(strings.TrimSpace(lastTag), "v")
	v := semver.New(version)
	v.BumpPatch()

	if !*getNext {
		reader := bufio.NewReader(os.Stdin)
		if _, err := fmt.Fprintf(os.Stderr, "Enter Release Version: [v%v] ", v); err != nil {
			panic(err)
		}

		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(text, "v") {
			text = text[1:]
			v = semver.New(strings.TrimSpace(text))
		}
		if _, err = fmt.Fprintf(os.Stderr, "Using Version: v%v\n", v); err != nil {
			panic(err)
		}
	}
	if *asNumber {
		fmt.Printf("%v", v)
	} else {
		fmt.Printf("v%v", v)
	}
}

func getLastTag() string {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		out = []byte("v0.0.0")
	}
	return string(out)
}

func getCurrentBranch() string {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}
