package obm

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	o := New("test")
	o.SetCap(100, 100)

	count := 100

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for i := 0; i < count; i++ {
		o.Asks.Put(Book{
			Price: r.NormFloat64() * 200,
			Size:  r.NormFloat64(),
		})
		o.Bids.Put(Book{
			Price: r.NormFloat64() * 100,
			Size:  r.NormFloat64(),
		})
	}

	// depth default 10
	books := o.Asks.Get(2)
	fmt.Println(books)
	fmt.Println("")
	sort.Sort(sort.Reverse(books))
	fmt.Println(books)

	sort.Sort(books)
	fmt.Println(books)
}
