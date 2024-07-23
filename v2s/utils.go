package obm

import "github.com/shopspring/decimal"

func (p *Orderbook) Remover(setPrice decimal.Decimal) *Orderbook {
	p.Lock()
	defer p.Unlock()

	p.Asks.tree.Each(func(key, val any) {
		s, ok := key.(string)
		if !ok {
			return
		}
		price, err := decimal.NewFromString(s)
		if err != nil {
			return
		}

		if setPrice.GreaterThan(price) {
			p.Asks.tree.Remove(key)
		}
	})
	p.Bids.tree.Each(func(key, val any) {
		s, ok := key.(string)
		if !ok {
			return
		}
		price, err := decimal.NewFromString(s)
		if err != nil {
			return
		}

		if setPrice.LessThan(price) {
			p.Bids.tree.Remove(key)
		}
	})

	return p
}
