package obm_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/go-numb/go-obm/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type Temp struct {
	MidPrice float64    `json:"mid_price"`
	Bids     []obm.Book `json:"bids"`
	Asks     []obm.Book `json:"asks"`
}

func TestLOB(t *testing.T) {
	o := obm.New("test")

	url := "https://api.bitflyer.com/v1/getboard?product_code=FX_BTC_JPY"
	res, err := http.Get(url)
	assert.NoError(t, err)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	temp := new(Temp)
	err = json.Unmarshal(b, temp)
	assert.NoError(t, err)

	fmt.Printf("ask input: %d\nbid input: %d\n", len(temp.Asks), len(temp.Bids))
	l := len(temp.Bids)
	if len(temp.Bids) < len(temp.Asks) {
		l = len(temp.Asks)
	}
	l *= 2
	o.SetCap(l, l)

	asks := make([]obm.Book, len(temp.Asks))
	copy(asks, temp.Asks)

	bids := make([]obm.Book, len(temp.Bids))
	copy(bids, temp.Bids)

	// inputが多く、capが少ない場合はMidPriceから昇順にして入力する
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price.LessThan(bids[j].Price)
	})

	o.Update(asks, bids)

	a, bb := o.Best()
	fmt.Printf("price: %s - size: %s\n", a.Price.String(), a.Size.String())
	fmt.Printf("price: %s - size: %s\n", bb.Price.String(), bb.Size.String())

	o.Asks.Each(func(key, val string) {
		dec, err := decimal.NewFromString(val)
		if err != nil {
			fmt.Println(err)
		}

		if a.Price.Equal(dec) {
			fmt.Printf("ask - %s - %s\n", key, val)
		}
	})

	fmt.Println("")

	o.Bids.Each(func(key, val string) {
		dec, err := decimal.NewFromString(val)
		if err != nil {
			fmt.Println(err)
		}

		if bb.Price.Equal(dec) {
			fmt.Printf("bid - %s - %s\n", key, val)
		}
	})

	r := 0.01
	lob := o.LOB(decimal.NewFromFloat(r))

	for i := 0; i < len(lob.Asks); i++ {
		fmt.Printf("%s - %s - %s - %s %%\n", lob.Asks[i].Price.String(), lob.Asks[i].Size.String(), lob.Asks[i].AccSize.String(), lob.Asks[i].AccRatio.Mul(decimal.NewFromInt(100)).String())
	}

	best := lob.Bids[0].Price
	for i := 0; i < len(lob.Bids); i++ {
		fmt.Printf("%s - %s - %s - %s %%\n", lob.Bids[i].Price.String(), lob.Bids[i].Size.String(), lob.Bids[i].AccSize.String(), lob.Bids[i].AccRatio.Mul(decimal.NewFromInt(100)).String())
	}

	fmt.Printf("%s - %s\n", best.String(), lob.Accumulation.String())
}

func TestAll(t *testing.T) {
	now := time.Now()
	defer func() {
		// exec time: 0.228133 s
		fmt.Printf("exec time: %f s\n", time.Since(now).Seconds())
	}()

	o := obm.New("test")

	url := "https://api.bitflyer.com/v1/getboard?product_code=FX_BTC_JPY"
	res, err := http.Get(url)
	assert.NoError(t, err)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	temp := new(Temp)
	err = json.Unmarshal(b, temp)
	assert.NoError(t, err)

	l := len(temp.Bids)
	if len(temp.Bids) < len(temp.Asks) {
		l = len(temp.Asks)
	}
	l *= 2
	o.SetCap(l, l)

	asks := make([]obm.Book, len(temp.Asks))
	copy(asks, temp.Asks)
	bids := make([]obm.Book, len(temp.Bids))
	copy(bids, temp.Bids)

	// inputが多く、capが少ない場合はMidPriceから昇順にして入力する
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price.LessThan(bids[j].Price)
	})

	o.Update(asks, bids)

	fmt.Printf("size: %#v\n", o.Asks.Size())
	o.Asks.Each(func(key, val string) {
		fmt.Printf("ask check: %s - %s\n", key, val)
	})
}
