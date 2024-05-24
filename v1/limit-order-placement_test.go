package obm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Temp struct {
	MidPrice float64 `json:"mid_price"`
	Bids     []Book  `json:"bids"`
	Asks     []Book  `json:"asks"`
}

func TestLOB(t *testing.T) {
	o := New("test")

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

	asks := make([]Book, len(temp.Asks))
	for i := range temp.Asks {
		asks[i] = temp.Asks[i]
	}
	bids := make([]Book, len(temp.Bids))
	for i := range temp.Bids {
		bids[i] = temp.Bids[i]
	}

	// inputが多く、capが少ない場合はMidPriceから昇順にして入力する
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price < bids[j].Price
	})

	o.Update(asks, bids)

	a, bb := o.Best()
	fmt.Printf("%f - %f\n", a.Price, a.Size)
	fmt.Printf("%f - %f\n", bb.Price, bb.Size)

	all := o.Asks._all()
	for i := 0; i < len(all); i++ {
		if i > 10 {
			break
		}
		if a.Price == all[i].Price {
			fmt.Printf("ask[%d] - %f - %f\n", i, all[i].Price, all[i].Size)
		}
	}
	fmt.Println("")

	o.Bids._all()
	all = o.Bids._all()
	for i := 0; i < len(all); i++ {
		if bb.Price == all[i].Price {
			fmt.Printf("bid[%d] - %f - %f\n", i, all[i].Price, all[i].Size)
		}
	}

	r := 0.01
	lob := o.LOB(r)

	for i := 0; i < len(lob.Asks); i++ {
		fmt.Printf("%.0f - %f - %f - %f%%\n", lob.Asks[i].Price, lob.Asks[i].Size, lob.Asks[i].AccSize, lob.Asks[i].AccRatio*100)
	}

	best := lob.Bids[0].Price
	for i := 0; i < len(lob.Bids); i++ {
		fmt.Printf("%.0f - %f - %f - %f - %f%%\n", lob.Bids[i].Price, 1-(lob.Bids[i].Price/best), lob.Bids[i].Size, lob.Bids[i].AccSize, lob.Bids[i].AccRatio*100)
	}

	fmt.Printf("%f - %f\n", best, lob.Accumulation)
}

func TestAll(t *testing.T) {
	now := time.Now()
	defer func() {
		// exec time: 0.197374 s
		fmt.Printf("exec time: %f s\n", time.Since(now).Seconds())
	}()

	o := New("test")

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

	asks := make([]Book, len(temp.Asks))
	for i := range temp.Asks {
		asks[i] = temp.Asks[i]
	}
	bids := make([]Book, len(temp.Bids))
	for i := range temp.Bids {
		bids[i] = temp.Bids[i]
	}

	// inputが多く、capが少ない場合はMidPriceから昇順にして入力する
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price < bids[j].Price
	})

	o.Update(asks, bids)

	a := o.Asks._all()
	fmt.Printf("size: %#v\n", o.Asks.Size())
	for i := 0; i < len(a); i++ {
		fmt.Printf("%.0f\n", a[i].Price)
	}
}
