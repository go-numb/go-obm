# Orderbook manager with Golang
go-obm is orderbook manager for float64. set, update, sort, get best price and more.  
go-obm/v2 supported decimal, key is string strict price.  

go-obm/v2s supported string keys and custom sort, like a v2.

## Installation

```
$ go get -u github.com/go-numb/go-obm

# v2 decimal
$ go get -u github.com/go-numb/go-obm/v2
```

## Usage
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	now := time.Now()
	defer func() {
	// # reference
	// ## v2 decimal exec time: 1.618996 s
	// ## v1 float64 exec time: 0.302507 s
		fmt.Printf("exec time: %f s\n", time.Since(now).Seconds())
	}()

	o := obm.New("test").SetCap(100000, 100000)

	count := 100000

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for i := 0; i < count; i++ {
		o.Asks.Put(obm.Book{
			Price: decimal.NewFromFloat(r.NormFloat64() * 200),
			Size:  decimal.NewFromFloat(r.NormFloat64()),
		})
		o.Bids.Put(obm.Book{
			Price: decimal.NewFromFloat(r.NormFloat64() * 100),
			Size:  decimal.NewFromFloat(r.NormFloat64()),
		})
	}

	// depth default 10
	books := o.Asks.Get(2)
	fmt.Println(books)
	fmt.Println("")
	sort.Sort(sort.Reverse(books))
	fmt.Println("reverse sort")
	fmt.Println(books)

	fmt.Println("sort")
	sort.Sort(books)
	fmt.Println(books)
}
// Print out
// 0 - -0.001213774853359162:-0.9522491064087244
// 1 - -0.00849450337623292:0.5419339371697502
// 2 - -0.012667106317365562:-0.273163161491901
// 3 - -0.024032114560307294:0.12993832989793802
// 4 - -0.02826975357672623:1.1198441111617725
// 5 - -0.035862938872954775:-0.6413617919818513
// 6 - -0.04878297355579431:-0.04172143152529045
// 7 - -0.04891801841933052:0.5753643028690486
// 8 - -0.05535350612112522:0.23952017070860265
// 9 - -0.05663744202366594:1.511194970209686
// reverse sort
// 0 - -0.05663744202366594:1.511194970209686
// 1 - -0.05535350612112522:0.23952017070860265
// 2 - -0.04891801841933052:0.5753643028690486
// 3 - -0.04878297355579431:-0.04172143152529045
// 4 - -0.035862938872954775:-0.6413617919818513
// 5 - -0.02826975357672623:1.1198441111617725
// 6 - -0.024032114560307294:0.12993832989793802
// 7 - -0.012667106317365562:-0.273163161491901
// 8 - -0.00849450337623292:0.5419339371697502
// 9 - -0.001213774853359162:-0.9522491064087244
// sort
// 0 - -0.001213774853359162:-0.9522491064087244
// 1 - -0.00849450337623292:0.5419339371697502
// 2 - -0.012667106317365562:-0.273163161491901
// 3 - -0.024032114560307294:0.12993832989793802
// 4 - -0.02826975357672623:1.1198441111617725
// 5 - -0.035862938872954775:-0.6413617919818513
// 6 - -0.04878297355579431:-0.04172143152529045
// 7 - -0.04891801841933052:0.5753643028690486
// 8 - -0.05535350612112522:0.23952017070860265
// 9 - -0.05663744202366594:1.511194970209686
// ## v2 decimal
// exec time: 1.618996 s

// ## v1 float64
// ...
// 7 - -822.442555:-0.193592
// 8 - -824.809016:-1.782220
// 9 - -833.094959:0.170305
// exec time: 0.302507 s
// PASS

	// limit order placement ratio of price range 
	r := 0.01
	lob := o.LOB(r)
	for i := 0; i < len(lob.Bids); i++ {
		fmt.Printf("%.0f - %f - %f - %f%%\n", lob.Bids[i].Price, lob.Bids[i].Size, lob.Bids[i].AccSize, lob.Bids[i].AccRatio*100)
	}

	// bid[996] - 5233998.000000 - 0.050000
	
	// price, size, accumulation of bids[A], ratio:[A]/accumulation of all
	// 5234465 - 0.010000 - 0.010000 - 0.031682%
	// 5234468 - 0.012767 - 0.022767 - 0.072130%
	// 5234522 - 0.012779 - 0.035546 - 0.112616%
	// 5234530 - 0.010000 - 0.045546 - 0.144298%
	// 5234539 - 0.015000 - 0.060546 - 0.191820%

}
```

## Author

[@_numbP](https://twitter.com/_numbP)

## License

[MIT](https://github.com/go-numb/go-obm/blob/master/LICENSE)