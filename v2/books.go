package obm

import (
	"fmt"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shopspring/decimal"
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
	return p.Books[j].Price.LessThan(p.Books[i].Price)
}

func (p *Books) Swap(i, j int) {
	p.Books[i], p.Books[j] = p.Books[j], p.Books[i]
}

type Book struct {
	Price decimal.Decimal
	Size  decimal.Decimal
}

func (p *Books) String() string {
	c := make([]Book, len(p.Books))
	copy(c, p.Books)

	s := make([]string, len(c))
	for i := range c {
		s[i] = fmt.Sprintf("%d - %s:%s", i, c[i].Price.String(), c[i].Size.String())
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
		b[i] = p.getBookValues(keys[i])

	}

	p.Books = b
	return p
}

func (p *Books) Put(book Book) {
	// put on when key is there
	if _, isThere := p.tree.Get(book.Price.String()); isThere {
		p.tree.Put(book.Price.String(), book.Size.String())
		return
	}

	// インプット情報はすべて入力保存
	//

	// If map size is the upper limit, delete data from the upper or lower limits.
	if p.tree.Size() >= p.cap {
		switch p.remover {
		case MAX:
			found, _ := p.tree.Max()
			price, err := decimal.NewFromString(found.(string))
			if err != nil {
				price = decimal.NewFromInt(0)
			}
			if book.Price.GreaterThan(price) {
				return
			}
			p.tree.Remove(found)

		case MIN:
			found, _ := p.tree.Min()
			price, err := decimal.NewFromString(found.(string))
			if err != nil {
				price = decimal.NewFromInt(0)
			}
			if book.Price.LessThan(price) {
				return
			}
			p.tree.Remove(found)

		}
	}

	p.tree.Put(book.Price.String(), book.Size.String())
}

func (p *Books) Each(fn func(key, val string)) {
	p.tree.Each(func(key, val any) {
		fn(key.(string), val.(string))
	})
}

func (p *Books) _all() []Book {
	var books []Book

	keys := p.tree.Keys()

	for i := range keys {
		books = append(books, p.getBookValues(keys[i]))
	}

	return books
}

func (p *Books) getBookValues(key any) Book {
	s := key.(string)
	price, err := decimal.NewFromString(s)
	if err != nil {
		price = decimal.Zero
	}
	val, _ := p.tree.Get(s)
	size, err := decimal.NewFromString(val.(string))
	if err != nil {
		size = decimal.Zero
	}
	return Book{
		Price: price,
		Size:  size,
	}
}

// Converter is a function that converts the price and size of the book.
func Converter(price, size any) Book {
	switch v := price.(type) {
	case int:
		return converterI(v, size.(int))
	case float32:
		return converterF32(v, size.(float32))
	case float64:
		return converterF(v, size.(float64))
	case string:
		return converterS(v, size.(string))

	case decimal.Decimal:
		if s, ok := size.(decimal.Decimal); ok {
			return Book{
				Price: v,
				Size:  s,
			}
		}
	}

	return Book{}
}

func converterI(price, size int) Book {
	return Book{
		Price: decimal.NewFromInt(int64(price)),
		Size:  decimal.NewFromInt(int64(size)),
	}
}

func converterF32(price, size float32) Book {
	return Book{
		Price: decimal.NewFromFloat32(price),
		Size:  decimal.NewFromFloat32(size),
	}
}

func converterF(price, size float64) Book {
	return Book{
		Price: decimal.NewFromFloat(price),
		Size:  decimal.NewFromFloat(size),
	}
}

func converterS(price, size string) Book {
	p, err := decimal.NewFromString(price)
	if err != nil {
		p = decimal.Zero
	}
	s, err := decimal.NewFromString(size)
	if err != nil {
		s = decimal.Zero
	}

	return Book{
		Price: p,
		Size:  s,
	}
}
