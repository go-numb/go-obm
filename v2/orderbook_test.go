package obm

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestOrderbook(t *testing.T) {
	now := time.Now()
	defer func() {
		// exec time: 0.036226 s
		// exec time: 58.981587 s
		fmt.Printf("exec time: %f s\n", time.Since(now).Seconds())
	}()

	o := New("BTC-PERP")
	o.SetCap(200, 200)

	l := 2000
	asks := make([]Book, l)
	bids := make([]Book, l)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := range asks {
		asks[i] = Book{
			Price: decimal.NewFromFloat(r.Float64()),
			Size:  decimal.NewFromFloat(r.NormFloat64() * 10),
		}
		time.Sleep(time.Nanosecond)
	}

	for i := range bids {
		bids[i] = Book{
			Price: decimal.NewFromFloat(r.Float64()),
			Size:  decimal.NewFromFloat(r.NormFloat64() * 10),
		}
		time.Sleep(time.Nanosecond)
	}

	start := time.Now()
	defer func() {
		fmt.Printf("make exec time: %f s\n", time.Since(start).Seconds())
	}()

	o.Update(asks, bids)

	fmt.Printf("updated time: %f s\n", time.Since(start).Seconds())

	fmt.Println(o.Asks.Get(10), "\n", o.Bids.Get(10))

	fmt.Printf("get time: %f s\n", time.Since(start).Seconds())

	fmt.Println(o.Best())

	fmt.Printf("ask: %d, bid: %d, %#v\n", len(o.Asks.Books), len(o.Bids.Books), o)

	fmt.Println(o.Wall())

}
