# Git hooks

Installation:

```shell
git config core.hooksPath develop/githooks
```

## pre-commit

The hook automatically runs `gofumpt`, `goimports`, and `goimports-reviser`
to correctly format the `.go` files included in the commit, provided
that all the aforementioned commands are installed and available
in the user's search `$PATH` environment variable:

```shell
go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/incu6us/goimports-reviser/v3@latest
```

It also runs `go mod tidy`.
