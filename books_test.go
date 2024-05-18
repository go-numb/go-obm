package obm

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	now := time.Now()
	defer func() {
		// exec time: 0.269598 s
		fmt.Printf("exec time: %f s\n", time.Since(now).Seconds())
	}()

	o := New("test").SetCap(100000, 100000)

	count := 100000

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
