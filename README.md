# go-talib

[![GoDoc](http://godoc.org/github.com/markcheno/go-talib?status.svg)](http://godoc.org/github.com/markcheno/go-talib) 

A pure [Go](http://golang.org/) port of [TA-Lib](http://ta-lib.org)

## Install

Install the package with:

```bash
go get github.com/markcheno/go-talib
```

Import it with:

```go
import "github.com/markcheno/go-talib"
```

and use `talib` as the package name inside the code.

## Example

```go
package main

import (
	"fmt"
	"github.com/markcheno/go-quote"
	"github.com/markcheno/go-talib"
)

func main() {
	spy, _ := quote.NewQuoteFromTiingo("spy", "2016-01-01", "2016-04-01", quote.Daily, true)
	fmt.Print(spy.CSV())
	rsi2 := talib.Rsi(spy.Close, 2)
	fmt.Println(rsi2)
}
```

Every indicator takes one or more `[]float64` price series (plus parameters) and
returns a `[]float64` (or a tuple of them) the same length as the input. The
leading "lookback" region that cannot be computed is zero-filled.

Inputs that are too short for the requested period, or an invalid (non-positive)
period, return a zero-filled slice of the correct length rather than panicking.

## Development

Tests compare every indicator against the reference Python
[TA-Lib](https://github.com/TA-Lib/ta-lib-python), so they need a `python` with
`numpy` and `TA-Lib` available. A [`justfile`](justfile) provisions that
interpreter with [uv](https://github.com/astral-sh/uv) automatically:

```bash
just test              # set up the venv (first run) and run the full suite
just test-one TestRsi  # run a single test or pattern
just bench MidP        # run benchmarks
```

`just setup` creates a project-local `.venv` (via uv) with `numpy` + `TA-Lib`;
the other recipes put it on `PATH` so the test harness's `python` resolves to
it. Run `just` with no arguments to list all recipes.

## License

MIT License  - see LICENSE for more details

# Contributors

- [Markcheno](https://github.com/markcheno) 
