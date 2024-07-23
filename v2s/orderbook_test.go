package obm

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOrderbook(t *testing.T) {
	start := time.Now()
	defer func() {
		// exec time: 0.036226 s
		// exec time: 58.981587 s
		fmt.Printf("exec time: %f s\n", time.Since(start).Seconds())
	}()

	var (
		setLength      = 100
		setUptedaCount = 10
		setMinFloat    = decimal.NewFromFloat(0.001)
		setMaxFloat    = decimal.NewFromFloat(1000)
	)

	o := New("BTC-PERP")
	o.SetCap(setLength, setLength)

	l := 50
	asks := make([]Book, l+1)
	bids := make([]Book, l+1)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	var (
		askmin = decimal.NewFromFloat(1000)
		askmax = decimal.NewFromFloat(-1)
		bidmin = decimal.NewFromFloat(1000)
		bidmax = decimal.NewFromFloat(-1)
	)

	for uptedaCount := 0; uptedaCount < setUptedaCount; uptedaCount++ {
		for i := range asks {
			switch i {
			case 0:
				asks[i] = Book{
					Price: setMinFloat,
					Size:  setMinFloat,
				}

			case 1:
				asks[i] = Book{
					Price: setMaxFloat,
					Size:  setMaxFloat,
				}

			default:
				asks[i] = Converter(r.Float64()+2, r.NormFloat64()+200)
			}

			if asks[i].Price.LessThan(askmin) {
				askmin = asks[i].Price
			}
			if asks[i].Price.GreaterThan(askmax) {
				askmax = asks[i].Price
			}

			time.Sleep(time.Nanosecond)
		}

		for i := range bids {
			switch i {
			case 0:
				bids[i] = Book{
					Price: setMinFloat,
					Size:  setMinFloat,
				}

			case 1:
				bids[i] = Book{
					Price: setMaxFloat,
					Size:  setMaxFloat,
				}

			default:
				bids[i] = Converter(r.Float64()+1, r.NormFloat64()+100)
			}

			if bids[i].Price.LessThan(bidmin) {
				bidmin = bids[i].Price
			}
			if bids[i].Price.GreaterThan(bidmax) {
				bidmax = bids[i].Price
			}

			time.Sleep(time.Nanosecond)
		}

		fmt.Printf("make exec time: %f s\n", time.Since(start).Seconds())

		o.Update(asks, bids)
		assert.Equal(t, o.Asks.Len(), len(o.Asks.Get(200).Books))
		assert.Equal(t, o.Bids.Len(), len(o.Bids.Get(200).Books))
		if o.Asks.Len() > setLength {
			assert.Equal(t, o.Asks.Len(), setLength)
		}
		if o.Bids.Len() > setLength {
			assert.Equal(t, o.Bids.Len(), setLength)
		}
	}

	amax, bmax := o.GetMax()
	amin, bmin := o.GetMin()
	fmt.Println("max", amax, bmin)

	bask, bbid := o.Best()
	// 売り板のベスト価格と売り板の
	assert.Equal(t, bask.Price.InexactFloat64(), amin.Price.InexactFloat64())
	assert.Equal(t, bbid.Price.InexactFloat64(), bmax.Price.InexactFloat64())

	// fmt.Printf("price: %s, %s, updated time: %f s\n", bbid.Price.String(), bask.Price.String(), time.Since(start).Seconds())

	// fmt.Println(o.Asks.Get(20), "\n", o.Bids.Get(20))

	fmt.Printf("get time: %f s\n", time.Since(start).Seconds())

	fmt.Println(o.Best())

	fmt.Printf("ask: %d, bid: %d, %#v\n", len(o.Asks.Books), len(o.Bids.Books), o)

	targetSize := 1
	fmt.Println(o.Wall(float64(targetSize)))

	maxPrice, maxSize := o.Asks.tree.Max()
	fmt.Printf("max price, size: %v, %v", maxPrice, maxSize)
	minPrice, minSize := o.Asks.tree.Min()
	fmt.Printf("min price, size: %v, %v", minPrice, minSize)

	o.Asks.Each(func(key, val string) {
		fmt.Println(key, val)
	})

}

func TestBook(t *testing.T) {
	orderbooks := New("BTC-PERP")

	maxLength := 1000

	orderbooks.SetCap(maxLength, maxLength)

	count := 2000

	var (
		bid = make([]Book, count)
		ask = make([]Book, count)
	)

	for i := 0; i < count; i++ {
		n := rand.Float64()

		bid[i] = Book{
			Price: decimal.NewFromFloat(n + 1),
			Size:  decimal.NewFromFloat(n + 100),
		}

		ask[i] = Book{
			Price: decimal.NewFromFloat(n + 2),
			Size:  decimal.NewFromFloat(n + 100),
		}

		fmt.Println(n)
	}

	orderbooks.Update(ask, bid)

	fmt.Printf("best: ")
	fmt.Println(orderbooks.Best())
	fmt.Printf("wall: ")
	fmt.Println(orderbooks.Wall(100.1))

	assert.Equal(t, orderbooks.Asks.Len(), maxLength)
}

func TestUpdate(t *testing.T) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	start := time.Now()
	defer func() {
		// exec time: 29.787411 s
		fmt.Printf("exec time: %f s\n", time.Since(start).Seconds())
	}()

	maxLength := 10000
	orderbooks := New("BTC-PERP").SetCap(maxLength, maxLength)

	updateCount := 20
	updateLength := 100
	// Update
	for i := 0; i < updateCount; i++ {
		var (
			asks = make([]Book, updateLength)
			bids = make([]Book, updateLength)
		)
		for j := 0; j < updateLength; j++ {
			price, _ := decimal.NewFromString(fmt.Sprintf("%f", float64(j)+r.NormFloat64()))
			size, _ := decimal.NewFromString(fmt.Sprintf("%f", float64(j)+r.NormFloat64()))

			asks = append(asks, Book{
				Price: price,
				Size:  size,
			})

			price, _ = decimal.NewFromString(fmt.Sprintf("%f", float64(j)+r.NormFloat64()))
			size, _ = decimal.NewFromString(fmt.Sprintf("%f", float64(j)+r.NormFloat64()))
			bids = append(bids, Book{
				Price: price,
				Size:  size,
			})

			time.Sleep(time.Nanosecond)
		}

		askbest, bidbest := orderbooks.Best()
		mid := askbest.Price.Add(bidbest.Price).Div(decimal.NewFromInt(2))
		orderbooks.Remover(mid)
		orderbooks.Update(asks, bids)

		expected := 10
		if i == 0 {
			expected = orderbooks.Bids.Size()
		}
		fmt.Printf("lentgh: ask %d, bid %d\n", orderbooks.Asks.Len(), orderbooks.Bids.Len())
		orderbooks.Asks.Get(10)
		orderbooks.Bids.Get(10)
		assert.Equal(t, orderbooks.Asks.Len(), expected)
		assert.Equal(t, orderbooks.Bids.Len(), expected)

		fmt.Printf("size: ask %d, bid %d\n", orderbooks.Asks.Size(), orderbooks.Bids.Size())

		ask, bid := orderbooks.Best()
		min, _ := orderbooks.Asks.tree.Min()
		assert.Equal(t, ask.Price.String(), min)
		max, _ := orderbooks.Bids.tree.Max()
		assert.Equal(t, bid.Price.String(), max)
	}
}
