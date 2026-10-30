package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bakito/semver/version"
	"github.com/coreos/go-semver/semver"
)

const (
	branchMain   = "main"
	branchMaster = "master"
	defaultTag   = "v0.0.0"
)

func formatVersion(v *semver.Version, numeric, withV bool) string {
	if !numeric && withV {
		return "v" + v.String()
	}
	return v.String()
}

func mustBeOnReleaseBranch(isNext, isCurrent bool, branch string) {
	if !isNext && !isCurrent && (branch != branchMain && branch != branchMaster) {
		panic(fmt.Errorf(`error: must be in "master/main" branch, current branch: %q`, branch))
	}
}

func promptVersion(defaultV *semver.Version, prefix string, withV bool) *semver.Version {
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fprintf(os.Stderr, "Enter Release Version: [%s] ", formatVersion(defaultV, false, withV)); err != nil {
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
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimPrefix(text, "v")
	return semver.New(strings.TrimSpace(text))
}

func main() {
	flagVersion := flag.Bool("version", false, "Print semver version and exit")
	flagNext := flag.Bool("next", false, "Just print the next version")
	flagCurrent := flag.Bool("current", false, "Just print the current version")
	flagNumeric := flag.Bool("numeric", false, "Numeric form")
	flagPrefix := flag.String("prefix", "", "Only consider tags with this prefix")
	flag.Parse()

	if *flagVersion {
		fmt.Printf("semver has version %s\n", version.Version)
		return
	}

	prefix, err := normalizePrefix(*flagPrefix)
	if err != nil {
		panic(err)
	}

	branch := getCurrentBranch()
	mustBeOnReleaseBranch(*flagNext, *flagCurrent, branch)

	v, withV, found := latestTag(getTags(prefix), prefix)
	if !found {
		v = semver.New(strings.TrimPrefix(defaultTag, "v"))
		withV = true
	}

	if *flagCurrent {
		fmt.Print(formatVersion(v, *flagNumeric, withV))
		return
	}

	v.BumpPatch()

	if *flagNext {
		fmt.Print(formatVersion(v, *flagNumeric, withV))
		return
	}

	v = promptVersion(v, prefix, withV)
	if _, err := fmt.Fprintf(os.Stderr, "Using Version: %s\n", formatVersion(v, false, withV)); err != nil {
		panic(err)
	}
	fmt.Print(formatVersion(v, *flagNumeric, withV))
}

func normalizePrefix(prefix string) (string, error) {
	if prefix == "" {
		return "", nil
	}
	if strings.HasPrefix(prefix, "-") {
		return "", fmt.Errorf("error: invalid prefix %q: must not start with '-'", prefix)
	}
	if strings.HasSuffix(prefix, "-") {
		return "", fmt.Errorf("error: invalid prefix %q: must not end with '-', it is added automatically", prefix)
	}
	return prefix + "-", nil
}

func getTags(prefix string) []string {
	if strings.HasPrefix(prefix, "-") {
		panic(fmt.Errorf("error: invalid prefix %q: must not start with '-'", prefix))
	}
	out, err := exec.Command("git", "tag", "--list", prefix+"*", "--merged", "HEAD").Output() //nolint:gosec // prefix is validated above and passed as a single argument without a shell
	if err != nil {
		return nil
	}
	return strings.Fields(string(out))
}

func latestTag(tags []string, prefix string) (*semver.Version, bool, bool) {
	var latest *semver.Version
	var latestWithV bool
	for _, tag := range tags {
		version, foundPrefix := strings.CutPrefix(tag, prefix)
		if !foundPrefix {
			continue
		}
		version, withV := strings.CutPrefix(version, "v")
		parsed, err := semver.NewVersion(version)
		if err != nil {
			continue
		}
		if latest == nil || parsed.Compare(*latest) > 0 {
			latest = parsed
			latestWithV = withV
		}
	}
	return latest, latestWithV, latest != nil
}

func getCurrentBranch() string {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}
