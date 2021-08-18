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