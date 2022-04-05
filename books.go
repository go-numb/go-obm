package obm

import (
	"fmt"
	"sort"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
)

type Books struct {
	cap     int
	remover Remover
	tree    *treemap.Map
	Books   []Book
}

type Book struct {
	Price float64
	Size  float64
}

func (p *Books) String() string {
	c := make([]Book, len(p.Books))
	copy(c, p.Books)

	s := make([]string, len(c))
	for i := range c {
		s[i] = fmt.Sprintf("%d - %f:%f", i, c[i].Price, c[i].Size)
	}
	return strings.Join(s, "\n")
}

// Get depth default:10
func (p *Books) Get(depth int) *Books {
	l := p.tree.Size()

	if depth < 0 {
		depth = 10
	}
	if depth > l {
		depth = l
	}

	b := make([]Book, depth)

	keys := p.tree.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].(float64) > keys[j].(float64)
	})

	for i := 0; i < depth; i++ {
		if value, isThere := p.tree.Get(keys[i]); isThere {
			b[i] = Book{
				Price: keys[i].(float64),
				Size:  value.(float64),
			}
		}
	}

	p.Books = b
	return p
}

func (p *Books) Put(book Book) {
	// put on when key is there
	if _, isThere := p.tree.Get(book.Price); isThere {
		p.tree.Put(book.Price, book.Size)
		return
	}

	// If map size is the upper limit, delete data from the upper or lower limits.
	if p.tree.Size() >= p.cap {
		switch p.remover {
		case MAX:
			found, _ := p.tree.Max()
			if book.Price > found.(float64) {
				return
			}
			p.tree.Remove(found)

		case MIN:
			found, _ := p.tree.Min()
			if book.Price < found.(float64) {
				return
			}
			p.tree.Remove(found)
		}

	}
	p.tree.Put(book.Price, book.Size)
}

func (p *Books) Each(fn func(key, val float64)) {
	p.tree.Each(func(key, val any) {
		fn(key.(float64), val.(float64))
	})
}
