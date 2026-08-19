# Repository Guidelines

## Project Structure & Module Organization

BaiduPCS-Go is a Go 1.23 CLI using `urfave/cli` v1. Command registration, the interactive shell, and completion live in `main.go`; keep command actions thin and delegate behavior to `internal/pcscommand/`. Netdisk API calls belong in `baidupcs/`, while download, upload, and transfer engines live under `internal/pcsfunctions/` and build on `requester/`. Shared leaf utilities are in `pcsutil/`. Account state and initialization are handled by `pcsconfig/` and `pcsinit/`. Tests are colocated with packages as `*_test.go`; deeper design notes are in `docs/` and `CODEBUDDY.md`.

Dependencies should flow from CLI code toward command logic, API/transfer layers, and finally utilities. Avoid introducing upward package dependencies.

## Build, Test, and Development Commands

- `go build`: build the local `BaiduPCS-Go` executable.
- `go run . ls`: run one CLI command; `go run .` starts the interactive shell.
- `go test ./...`: run all unit tests without requiring a live Baidu account.
- `gofmt -w <files>`: format changed Go files.
- `go vet ./...`: run the repository's static checks.
- `./build.sh vX.Y.Z`: produce cross-platform release archives with an injected version.

Set `BAIDUPCS_GO_VERBOSE=1` or pass `--verbose` when debugging requests.

## Coding Style & Naming Conventions

Follow standard Go formatting and naming: tabs via `gofmt`, short lowercase package names, exported identifiers in `PascalCase`, and internal identifiers in `camelCase`. Prefer one command group per file in `internal/pcscommand/`. Return typed errors through `baidupcs/pcserror`; do not add library `panic`, `unsafe`, or `//go:linkname` usage. Preserve `urfave/cli` v1 APIs.

## Testing Guidelines

Use Go's `testing` package and name tests `TestBehavior`. Add focused tests beside modified code, especially for request signing, argument parsing, checksums, and transfer edge cases. There is no stated coverage threshold, but bug fixes should include regression tests. Do not add tests that require real credentials or network access.

## Commit & Pull Request Guidelines

History favors concise prefixes such as `fix:`, `feat:`, `chore:`, and `Refactor`. Prefer an imperative English subject, optionally scoped: `fix(transfer): handle public shares`. Pull requests should explain the problem and solution, link relevant issues, list validation commands, and note CLI-visible or configuration changes. Target the appropriate long-lived branch (`main` for releases or `feature/dev` for development).

## Security & Configuration Notes

Never commit BDUSS, STOKEN, account files, or captured request data. Do not casually change app IDs, user-agent constants, signing code in `baidupcs/netdisksign/`, or non-SVIP concurrency defaults; these are compatibility-sensitive. Treat the global mutable `pcsconfig.Config` carefully in concurrent code and tests.
