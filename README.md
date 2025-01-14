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

## License

MIT License  - see LICENSE for more details

# Contributors

- [Markcheno](https://github.com/markcheno) 
