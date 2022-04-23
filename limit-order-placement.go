package obm

import "sort"

type LimitOrderPlacement struct {
	Accumulation float64
	AccAsksRatio float64
	AccBidsRatio float64

	Asks []LimitOrderBook
	Bids []LimitOrderBook
}

type LimitOrderBook struct {
	Book
	AccSize  float64
	AccRatio float64
}

// LOB 範囲板(cap)の総量とOrderPlacePriceの累計量と割合を算出
func (p *Orderbook) LOB(rangeRatio float64) *LimitOrderPlacement {
	p.Lock()
	defer p.Unlock()

	ask, bid := p.Best()
	mid := (ask.Price + bid.Price) * 0.5
	maxPrice := mid * (1 + rangeRatio)
	minPrice := mid * (1 - rangeRatio)

	asks := p.Asks._all()
	bids := p.Bids._all()
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price > bids[j].Price
	})

	// sum
	var a, b float64
	for i := range asks {
		if maxPrice < asks[i].Price {
			break
		}
		a += asks[i].Size
	}
	for i := range bids {
		if minPrice > bids[i].Price {
			break
		}
		b += bids[i].Size
	}
	acc := a + b

	lob := new(LimitOrderPlacement)
	lob.Accumulation = acc

	// Askの累積は小さいPriceから
	// Bidの累積は大きいPriceから
	var (
		accA  float64
		lasks []LimitOrderBook
	)
	for i := range asks {
		// 指定価格「超過」の範囲外
		if maxPrice < asks[i].Price {
			break
		}

		accA += asks[i].Size
		lasks = append(lasks, LimitOrderBook{
			Book: Book{
				Price: asks[i].Price,
				Size:  asks[i].Size,
			},
			AccSize:  accA,
			AccRatio: accA / acc,
		})
	}
	lob.Asks = lasks
	lob.AccAsksRatio = accA / acc

	var (
		accB  float64
		lbids []LimitOrderBook
	)
	for i := range bids {
		// 指定価格「未満」の範囲外
		if minPrice > bids[i].Price {
			break
		}

		accB += bids[i].Size
		lbids = append(lbids, LimitOrderBook{
			Book: Book{
				Price: bids[i].Price,
				Size:  bids[i].Size,
			},
			AccSize:  accB,
			AccRatio: accB / acc,
		})
	}
	lob.Bids = lbids
	lob.AccBidsRatio = accB / acc

	return lob
}
