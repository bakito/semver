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

// Introduced constants for well-known values.
const (
	branchMain   = "main"
	branchMaster = "master"
	defaultTag   = "v0.0.0"
)

// Extracted: centralize version formatting to avoid duplication.
func formatVersion(v *semver.Version, numeric bool) string {
	if numeric {
		return v.String()
	}
	return "v" + v.String()
}

// Extracted: encapsulate release-branch validation for clarity.
func mustBeOnReleaseBranch(isNext, isCurrent bool, branch string) {
	if !isNext && !isCurrent && (branch != branchMain && branch != branchMaster) {
		panic(fmt.Errorf(`error: must be in "master/main" branch, current branch: %q`, branch))
	}
}

// Extracted: prompt the user for a version, with default shown and validation.
func promptVersion(defaultV *semver.Version) *semver.Version {
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fprintf(os.Stderr, "Enter Release Version: [v%v] ", defaultV); err != nil {
		panic(err)
	}
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultV
	}
	text = strings.TrimPrefix(text, "v")
	return semver.New(strings.TrimSpace(text))
}

// ... existing code ...
func main() {
	flagNext := flag.Bool("next", false, "Just print the next version")
	flagCurrent := flag.Bool("current", false, "Just print the current version")
	flagNumeric := flag.Bool("numeric", false, "Numeric form")
	flag.Parse()

	branch := getCurrentBranch()
	mustBeOnReleaseBranch(*flagNext, *flagCurrent, branch)

	lastTag := getLastTag()
	version := strings.TrimPrefix(strings.TrimSpace(lastTag), "v")
	v := semver.New(version)

	if *flagCurrent {
		fmt.Print(formatVersion(v, *flagNumeric))
		return
	}

	v.BumpPatch()

	if *flagNext {
		fmt.Print(formatVersion(v, *flagNumeric))
		return
	}

	v = promptVersion(v)
	if _, err := fmt.Fprintf(os.Stderr, "Using Version: %s\n", formatVersion(v, false)); err != nil {
		panic(err)
	}
}

func getLastTag() string {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		out = []byte(defaultTag)
	}
	return strings.TrimSpace(string(out))
}

func getCurrentBranch() string {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}
