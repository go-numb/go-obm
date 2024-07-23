package obm

import (
	"sort"

	"github.com/shopspring/decimal"
)

type LimitOrderPlacement struct {
	Accumulation decimal.Decimal
	AccAsksRatio decimal.Decimal
	AccBidsRatio decimal.Decimal

	Asks []LimitOrderBook
	Bids []LimitOrderBook
}

type LimitOrderBook struct {
	Book
	AccSize  decimal.Decimal
	AccRatio decimal.Decimal
}

// LOB 範囲板(cap)の総量とOrderPlacePriceの累計量と割合を算出
func (p *Orderbook) LOB(rangeRatio decimal.Decimal) *LimitOrderPlacement {
	p.Lock()
	defer p.Unlock()

	ask, bid := p.Best()
	mid := ask.Price.Add(bid.Price).Div(decimal.NewFromInt(2))
	maxPrice := mid.Mul(decimal.NewFromInt(1).Add(rangeRatio))
	minPrice := mid.Mul(decimal.NewFromInt(1).Sub(rangeRatio))

	asks := p.Asks._all()
	bids := p.Bids._all()
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price.GreaterThan(bids[j].Price)
	})

	var (
		a, b decimal.Decimal
	)
	for i := range asks {
		if maxPrice.LessThan(asks[i].Price) {
			break
		}
		a = a.Add(asks[i].Size)
	}
	for i := range bids {
		if minPrice.GreaterThan(bids[i].Price) {
			break
		}
		b = b.Add(bids[i].Size)
	}
	acc := a.Add(b)

	lob := new(LimitOrderPlacement)
	lob.Accumulation = acc

	var (
		accA  decimal.Decimal
		lasks []LimitOrderBook
	)
	for i := range asks {
		if maxPrice.LessThan(asks[i].Price) {
			break
		}

		accA = accA.Add(asks[i].Size)
		lasks = append(lasks, LimitOrderBook{
			Book: Book{
				Price: asks[i].Price,
				Size:  asks[i].Size,
			},
			AccSize:  accA,
			AccRatio: accA.Div(acc),
		})
	}
	lob.Asks = lasks
	lob.AccAsksRatio = accA.Div(acc)

	var (
		accB  decimal.Decimal
		lbids []LimitOrderBook
	)
	for i := range bids {
		if minPrice.GreaterThan(bids[i].Price) {
			break
		}

		accB = accB.Add(bids[i].Size)
		lbids = append(lbids, LimitOrderBook{
			Book: Book{
				Price: bids[i].Price,
				Size:  bids[i].Size,
			},
			AccSize:  accB,
			AccRatio: accB.Div(acc),
		})
	}
	lob.Bids = lbids
	lob.AccBidsRatio = accB.Div(acc)

	return lob
}
