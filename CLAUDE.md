# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`go-talib` is a pure Go port of [TA-Lib](http://ta-lib.org), the technical analysis library. It is a single-package library (`package talib`) with no third-party dependencies — everything lives in [talib.go](talib.go) (~5800 lines) with tests in [talib_test.go](talib_test.go).

## Commands

```bash
go build ./...        # build
go vet ./...          # vet
go test ./...         # run all tests (REQUIRES python + talib installed — see below)
go test -run TestRsi  # run a single indicator's test
```

## Testing model (important)

The tests do **not** hardcode expected values. Instead, each test computes a result in Go and shells out to the reference Python TA-Lib implementation, then compares the two arrays element-by-element. Consequences:

- `go test` requires a working `python` on PATH with the `talib` package importable. `TestMain` aborts the whole suite with a non-zero exit if `python -c "import talib"` fails. There is no way to run the suite without it.
- `compare(t, goResult, "result = talib.RSI(testClose,10)")` is the core helper: the second arg is a literal snippet of Python that must assign to a variable named `result`. The Python snippet runs against the same `testOpen/testHigh/testLow/testClose/testVolume/testRand` fixture arrays, which are injected as numpy arrays.
- Comparison is deliberately fuzzy: values are formatted to 6 decimals (reduced precision for magnitudes ≥ 1e6) and the **last two digits are dropped** before comparing, to tolerate floating-point divergence from the C/Python reference. NaN on either side maps to `0.0`.
- When adding/porting a function, add a `TestXxx` that calls `compare` with the equivalent `talib.XXX(...)` Python call — that *is* the correctness oracle.

## Architecture / conventions

The library mirrors TA-Lib's C API shape rather than being idiomatic Go:

- **Slice in, slice out.** Functions take one or more `[]float64` price series (plus int periods / float params) and return `[]float64` (or a tuple of slices for multi-output indicators like `BBands`, `Macd`, `Stoch`, `Aroon`). Output slices are always the **same length as the input**; the leading "lookback" region that TA-Lib can't compute is zero-filled rather than truncated. Callers index by position, so preserving this alignment is essential.
- **Naming.** Exported functions use Go-cased names of TA-Lib functions (`Rsi`, `Macd`, `BBands`, `HtTrendline`). Parameters keep the TA-Lib `inReal`, `inHigh`, `inTimePeriod`, `inMAType` naming. Each function has a doc comment with a Python-style signature line documenting defaults.
- **Moving-average dispatch.** `MaType` is an int enum (`SMA, EMA, WMA, DEMA, TEMA, TRIMA, KAMA, MAMA, T3MA`). `Ma(...)` switches on it; many indicators accept an `inMAType` and delegate through it. Unexported lowercase helpers (e.g. `ema(inReal, period, k1)`) hold shared math that the exported wrappers call with specific smoothing constants.
- **Sections.** [talib.go](talib.go) is organized by TA-Lib category, marked with banner comments: Overlap Studies, Momentum Indicators, Volume Indicators, Volatility Indicators, Price Transform, Cycle Indicators, Statistic Functions, Math Transform Functions, Math Operator Functions. Keep new functions in the matching section.

When porting a new indicator, follow the existing pattern literally: match TA-Lib's lookback/zero-fill behavior exactly, then prove it with a `compare`-based test against Python TA-Lib.
