package obm_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	obmv2s "github.com/go-numb/go-obm/v2s"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	now := time.Now()
	defer func() {
		// exec time: 2.114597 s
		fmt.Printf("exec time: %f s\n", time.Since(now).Seconds())
	}()

	count := 10000
	o := obmv2s.New("test").SetCap(count, count)

	prices := make([]float64, count)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for i := 0; i < count; i++ {
		price, _ := decimal.NewFromString(fmt.Sprintf("%f", r.NormFloat64()*200))
		size, _ := decimal.NewFromString(fmt.Sprintf("%f", r.NormFloat64()))

		prices = append(prices, price.InexactFloat64())
		o.Asks.Put(obmv2s.Book{
			Price: price,
			Size:  size,
		})

		price, _ = decimal.NewFromString(fmt.Sprintf("%f", r.NormFloat64()*100))
		size, _ = decimal.NewFromString(fmt.Sprintf("%f", r.NormFloat64()))
		o.Bids.Put(obmv2s.Book{
			Price: price,
			Size:  size,
		})
	}

	assert.Equal(t, count, o.Asks.Size())

	for i := 0; i < len(prices); i++ {
		fmt.Printf("array: %f\n", prices[i])
	}
	fmt.Println("")

	// depth default 10
	books := o.Asks.Get(2)
	fmt.Println(books)
	assert.Equal(t, 10, len(books.Books))
	fmt.Println("")
	sort.Sort(sort.Reverse(books))
	fmt.Println("reverse sort")
	fmt.Println(books)

	fmt.Println("sort")
	sort.Sort(books)
	fmt.Println(books)
}
