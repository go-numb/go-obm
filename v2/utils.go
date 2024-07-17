package obm

func (p *Orderbook) Remover(setPrice float64) *Orderbook {
	p.Lock()
	defer p.Unlock()

	p.Asks.tree.Each(func(key, val any) {
		price := key.(float64)
		if setPrice > price {
			p.Asks.tree.Remove(key)
		}
	})
	p.Bids.tree.Each(func(key, val any) {
		price := key.(float64)
		if setPrice < price {
			p.Bids.tree.Remove(key)
		}
	})

	return p
}
