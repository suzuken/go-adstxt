// Package adstxt implements Ads.txt protocol defined by iab.
package adstxt

import (
	"io"
	"net/http"
)

// Record is ads.txt data field defined in iab.
type Record struct {
	// ExchangeDomain is domain name of the advertising system
	ExchangeDomain string

	// PublisherAccountID is the identifier associated with the seller
	// or reseller account within the advertising system.
	PublisherAccountID string

	// AccountType is an enumeration of the type of account.
	AccountType AccountType

	// AuthorityID is an ID that uniquely identifies the advertising system
	// within a certification authority.
	AuthorityID string
}

const (
	AccountDirect AccountType = iota
	AccountReseller
	AccountOther
)

type AccountType int

func Get(rawurl string) ([]Record, error) {
	resp, err := http.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return Parse(resp.Body)
}

func Parse(in io.Reader) ([]Record, error) {
	records := make([]Record, 0)
	p := NewParser(in)

LOOP:
	for {
		r, err := p.Parse()
		if err == io.EOF {
			break LOOP
		}
		if err != nil {
			return nil, err
		}
		if r != nil {
			records = append(records, *r)
		}
	}

	return records, nil
}
