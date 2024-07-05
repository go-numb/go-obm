package obm

import (
	"sort"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/shopspring/decimal"
)

type Remover int

const (
	MAX Remover = iota
	MIN
)

type Orderbook struct {
	sync.Mutex

	Symbol    string
	Bids      *Books
	Asks      *Books
	UpdatedAt time.Time
}

// New is Create a new Orderbook
// default cap is 5
func New(symbol string) *Orderbook {
	return &Orderbook{
		Symbol: symbol,
		Asks: &Books{
			isBid: false,
			cap:   5,
			// ascending order
			remover: MAX,
			tree:    treemap.NewWith(utils.Float64Comparator),
			Books:   []Book{},
		},
		Bids: &Books{
			isBid: true,
			cap:   5,
			// descending-order
			remover: MIN,
			tree:    treemap.NewWith(utils.Float64Comparator),
			Books:   []Book{},
		},
		UpdatedAt: time.Now(),
	}
}

// SetCap is Determine the upper and lower limits of length stored in Map
func (p *Orderbook) SetCap(askcap, bidcap int) *Orderbook {
	p.Bids.cap = bidcap
	p.Asks.cap = askcap

	return p
}

func (p *Orderbook) GetCap() (askcap, bidcap int) {
	return p.Asks.cap, p.Bids.cap
}

func (p *Orderbook) GetMin() (askmin, bidmin Book) {
	pa, sa := p.Asks.tree.Min()
	pb, sb := p.Bids.tree.Min()
	return Converter(pa, sa), Converter(pb, sb)
}

func (p *Orderbook) GetMax() (askmax, bidmax Book) {
	pa, sa := p.Asks.tree.Max()
	pb, sb := p.Bids.tree.Max()
	return Converter(pa, sa), Converter(pb, sb)
}

func (p *Orderbook) Update(asks, bids []Book) {
	p.Lock()
	defer p.Unlock()

	// Which is faster, Workgroup or exec bids to asks, depends on the array length.
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range asks {
			if asks[i].Size.LessThanOrEqual(decimal.Zero) {
				p.Asks.tree.Remove(asks[i].Price.InexactFloat64())
			} else {
				p.Asks.Put(asks[i])
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range bids {
			if bids[i].Size.LessThanOrEqual(decimal.Zero) {
				p.Bids.tree.Remove(bids[i].Price.InexactFloat64())
			} else {
				p.Bids.Put(bids[i])
			}
		}
	}()

	wg.Wait()

	p.UpdatedAt = time.Now()
}

func (p *Orderbook) Best() (ask, bid Book) {
	aprice, asize := p.Asks.tree.Min()
	bprice, bsize := p.Bids.tree.Max()

	if aprice == nil || asize == nil || bprice == nil || bsize == nil {
		return Book{}, Book{}
	}

	// Convert aprice and bprice from string to decimal.Decimal
	a := Converter(aprice, asize)
	b := Converter(bprice, bsize)

	return a, b
}

// Wall search Big board In the setting cap range
// Search by Price near Mid
func (p *Orderbook) Wall() (ask, bid Book) {
	p.Lock()
	defer p.Unlock()

	wg := sync.WaitGroup{}

	// Ask is descending-order
	wg.Add(1)
	go func() {
		defer wg.Done()

		p.Asks.Each(func(key, val float64) {
			b := Converter(key, val)

			if ask.Size.LessThan(b.Size) {
				ask.Price = b.Price
				ask.Size = b.Size
			}
		})
	}()

	// Bid is ascending-order
	wg.Add(1)
	go func() {
		defer wg.Done()

		prices := p.Bids.tree.Keys()
		sort.Slice(prices, func(i, j int) bool {
			if prices[i] == nil {
				p.Bids.tree.Remove(prices[i])
				return false
			}
			if prices[j] == nil {
				p.Bids.tree.Remove(prices[j])
				return false
			}

			a := decimal.NewFromFloat(prices[i].(float64))
			b := decimal.NewFromFloat(prices[j].(float64))

			return a.GreaterThan(b)
		})

		for i := 0; i < len(prices); i++ {
			if v, isThere := p.Bids.tree.Get(prices[i]); isThere {
				b := Converter(prices[i], v)

				if bid.Size.LessThan(b.Size) {
					bid.Price = b.Price
					bid.Size = b.Size
				}
			}
		}
	}()

	wg.Wait()

	return
}
