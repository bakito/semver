[![Go](https://github.com/bakito/semver/actions/workflows/go.yml/badge.svg)](https://github.com/bakito/semver/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bakito/semver)](https://goreportcard.com/report/github.com/bakito/semver)
[![GitHub Release](https://img.shields.io/github/release/bakito/semver.svg?style=flat)](https://github.com/bakito/semver/releases)
# semver

semver allows to evaluate the current release tag of a git repository.
It may help releasing application released with [goreleaser](https://github.com/goreleaser/goreleaser).

## Install 

```bash
go get -u github.com/bakito/semver
```

## Usage
```bash
version=$(semver); \
git tag -s $version -m"Release $version"
goreleaser --rm-dist
```