package adstxt_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/suzuken/go-adstxt"
)

func TestParseAdstxt(t *testing.T) {
	cases := []struct {
		txt      string
		expected []adstxt.Record
	}{
		{
			txt: `example.com,1,DIRECT`,
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
				},
			},
		},
		{
			txt: "example.com,1,DIRECT\nexample.org,2,RESELLER",
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
				},
				{
					ExchangeDomain:     "example.org",
					PublisherAccountID: "2",
					AccountType:        adstxt.AccountReseller,
				},
			},
		},
		{
			txt: "# comment out\nexample.com,1,DIRECT",
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
				},
			},
		},
	}

	for i, c := range cases {
		record, err := adstxt.Parse(strings.NewReader(c.txt))
		if err != nil {
			t.Errorf("(#%d) parse ads.txt failed: %s", i, err)
		}
		if !reflect.DeepEqual(c.expected, record) {
			t.Errorf("want %v, got %v", c.expected, record)
		}
	}
}
