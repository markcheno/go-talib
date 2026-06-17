# go-talib tasks.
#
# The test suite (talib_test.go) shells out to a `python` that must have numpy
# and the reference TA-Lib installed, and compares its output against this Go
# port. These recipes provision that interpreter with uv and put it on PATH.

# Locate uv (it isn't always on PATH); fall back to the standard install dir.
uv := `command -v uv 2>/dev/null || echo "$HOME/.local/bin/uv"`

# Project-local virtualenv holding numpy + TA-Lib.
venv := justfile_directory() / ".venv"
py   := venv / "bin" / "python"

# List available recipes.
default:
    @just --list

# Create the uv venv and install numpy + TA-Lib (idempotent; uv caches).
setup:
    {{uv}} venv {{venv}}
    {{uv}} pip install --python {{py}} numpy TA-Lib

# Ensure the venv exists and has numpy + talib importable, else (re)create it.
_ensure-venv:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! "{{py}}" -c 'import talib, numpy' >/dev/null 2>&1; then
        just setup
    fi

# Run the full parity test suite.
test: _ensure-venv
    PATH="{{venv}}/bin:$PATH" go test -count=1 ./...

# Run a single test or pattern, e.g. `just test-one TestRsi`.
test-one pattern: _ensure-venv
    PATH="{{venv}}/bin:$PATH" go test -count=1 -run '{{pattern}}' -v ./...

# Run benchmarks (TestMain still gates the package, so the venv is needed).
bench pattern=".": _ensure-venv
    PATH="{{venv}}/bin:$PATH" go test -run '^$' -bench '{{pattern}}' -benchmem ./...

# Build and vet (no Python needed).
build:
    go build ./...

vet:
    go vet ./...

# Format check.
fmt:
    gofmt -l .

# Remove the venv.
clean:
    rm -rf {{venv}}
