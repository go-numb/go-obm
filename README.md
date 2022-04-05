# Orderbook manager with Golang
go-obm is orderbook manager. set, update, sort, get best price and more.

## Installation

```
$ go get -u github.com/go-numb/go-obm
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
    // Setup
	o := New("BTC-PERP")
	o.SetCap(20, 20)

    // input(server response...)
    // dummy data
	l := 10
	asks := make([]Book, l)
	bids := make([]Book, l)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := range asks {
		asks[i] = Book{
			Price: r.Float64(),
			Size:  r.NormFloat64() * 10,
		}
		time.Sleep(time.Millisecond)
	}

	for i := range bids {
		bids[i] = Book{
			Price: r.Float64(),
			Size:  r.NormFloat64() * 10,
		}
		time.Sleep(time.Millisecond)
	}

	o.Update(asks, bids)


    // use struct
	fmt.Println(o.Asks.Get(10), "\n", o.Bids.Get(10))

	fmt.Println(o.Best())

	fmt.Printf("ask: %d, bid: %d, %#v\n", len(o.Asks.Books), len(o.Bids.Books), o)

// Print out

// updated time: 0.000000 s

// asks------------------------
// 0 - 0.999126:12.623239
// 1 - 0.998673:11.969490
// 2 - 0.995374:12.411968
// 3 - 0.993730:19.209979
// 4 - 0.993432:1.188730
// 5 - 0.993338:10.510313
// 6 - 0.993020:17.003518
// 7 - 0.988681:10.699286
// 8 - 0.986657:6.555007
// 9 - 0.983770:17.040381 
// bids------------------------
// 0 - 0.158852:2.346508
// 1 - 0.157248:4.807470
// 2 - 0.154350:7.245589
// 3 - 0.150005:2.467844
// 4 - 0.148196:1.627250
// 5 - 0.146792:13.557464
// 6 - 0.144692:4.510992
// 7 - 0.134729:2.586043
// 8 - 0.133824:2.608753
// 9 - 0.132424:3.714124

// get time: 0.000000 s

// bestask{0.765535581376781 8.375973510078225}, bestbid{0.15885186778959806 2.346508140219099}
// ask: 10, bid: 10, &obm.Orderbook{Mutex:sync.Mutex{state:0, sema:0x0}, Symbol:"BTC-PERP", Bids:(*obm.Books)(0xc000074750), Asks:(*obm.Books)(0xc000074780), UpdatedAt:time.Date(2022, time.April, 5, 20, 15, 59, 532576600, time.Local)}
// exec time: 0.009937 s
// PASS

}
```

## Author

[@_numbP](https://twitter.com/_numbP)

## License

[MIT](https://github.com/go-numb/go-obm/blob/master/LICENSE)