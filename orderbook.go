package obm

import (
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
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

func New(symbol string) *Orderbook {
	return &Orderbook{
		Symbol: symbol,
		Bids: &Books{
			cap: 0,
			// descending-order
			remover: MAX,
			tree:    treemap.NewWith(utils.Float64Comparator),
			Books:   []Book{},
		},
		Asks: &Books{
			cap: 0,
			// ascending order
			remover: MIN,
			tree:    treemap.NewWith(utils.Float64Comparator),
			Books:   []Book{},
		},
		UpdatedAt: time.Now(),
	}
}

// SetCap is Determine the upper and lower limits of length stored in Map
func (p *Orderbook) SetCap(bidcap, askcap int) {
	p.Bids.cap = bidcap
	p.Asks.cap = askcap
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
			if asks[i].Size <= 0 {
				p.Asks.tree.Remove(asks[i].Price)
			} else {
				p.Asks.Put(asks[i])
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range bids {
			if bids[i].Size <= 0 {
				p.Bids.tree.Remove(bids[i].Price)
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
	return Book{aprice.(float64), asize.(float64)}, Book{bprice.(float64), bsize.(float64)}
}
