// Package adstxt implements Ads.txt protocol defined by iab.
package adstxt

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func parseAccountType(s string) AccountType {
	switch s {
	case "Direct", "DIRECT":
		return AccountDirect
	case "Reseller", "RESELLER":
		return AccountReseller
	default:
		// NOTE or should be error ?
		return AccountOther
	}
}

func parseRow(row string) (Record, error) {
	fields := strings.Split(row, ",")
	if len := len(fields); len != 3 && len != 4 {
		return Record{}, fmt.Errorf("ads.txt has fields length is incorrect.: %s", row)
	}

	var r Record
	r.ExchangeDomain = fields[0]
	r.PublisherAccountID = fields[1]
	r.AccountType = parseAccountType(fields[2])

	// AuthorityID is optional
	if len(fields) >= 4 {
		r.AuthorityID = fields[3]
	}
	return r, nil
}

func Parse(in io.Reader) ([]Record, error) {
	records := make([]Record, 0)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		text := scanner.Text()

		// comment out
		if []rune(text)[0] == '#' || strings.Index(text, "contact=") == 0 {
			continue
		}
		r, err := parseRow(scanner.Text())
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
