package obm

import (
	"fmt"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
)

type Books struct {
	cap     int
	remover Remover
	tree    *treemap.Map
	Books   []Book
}

func (p *Books) Len() int {
	return len(p.Books)
}

func (p *Books) Less(i, j int) bool {
	return p.Books[i].Price > p.Books[j].Price
}

func (p *Books) Swap(i, j int) {
	p.Books[i], p.Books[j] = p.Books[j], p.Books[i]
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

func (p *Books) Size() int {
	return p.tree.Size()
}

// Get depth default:10
func (p *Books) Get(depth int) *Books {
	l := p.tree.Size()

	if depth < 10 {
		depth = 10
	}
	if depth > l {
		depth = l
	}

	b := make([]Book, depth)

	// sorted
	keys := p.tree.Keys()

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

	// インプット情報はすべて入力保存
	//

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

func (p *Books) _all() []Book {
	var books []Book

	keys := p.tree.Keys()

	for i := range keys {
		if value, isThere := p.tree.Get(keys[i]); isThere {
			books = append(books, Book{
				Price: keys[i].(float64),
				Size:  value.(float64),
			})
		}
	}

	return books
}
