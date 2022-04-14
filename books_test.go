package obm_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/go-numb/go-obm"
)

func TestSort(t *testing.T) {
	o := obm.New("test")
	o.SetCap(100, 100)

	count := 100

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for i := 0; i < count; i++ {
		o.Asks.Put(obm.Book{
			Price: r.NormFloat64() * 200,
			Size:  r.NormFloat64(),
		})
		o.Bids.Put(obm.Book{
			Price: r.NormFloat64() * 100,
			Size:  r.NormFloat64(),
		})
	}

	// depth default 10
	books := o.Asks.Get(2)
	fmt.Println(books)

	sort.Sort(sort.Reverse(books))
	fmt.Println(books)
}
