package obm_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/go-numb/go-obm/v2"
	"github.com/shopspring/decimal"
)

func TestSort(t *testing.T) {
	now := time.Now()
	defer func() {
		// exec time: 2.114597 s
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
