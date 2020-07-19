// Package adstxt implements Ads.txt protocol defined by iab.
package adstxt

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
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
	switch strings.ToUpper(s) {
	case "DIRECT":
		return AccountDirect
	case "RESELLER":
		return AccountReseller
	default:
		// NOTE or should be error ?
		return AccountOther
	}
}

var leadingBlankRe = regexp.MustCompile(`\A[\s\t]+`)
var trailingBlankRe = regexp.MustCompile(`[\s\t]+\z`)

func normalize(s string) string {
	// sanitize blank characters
	s = leadingBlankRe.ReplaceAllString(s, "")
	s = trailingBlankRe.ReplaceAllString(s, "")
	return s
}

func parseRow(row string) (*Record, error) {
	// dropping extension field
	if idx := strings.Index(row, ";"); idx != -1 {
		row = row[0:idx]
	}

	fields := strings.Split(row, ",")

	// if the first field contains "=", then the row is for key-value definitions
	if strings.Index(fields[0], "=") != -1 {
		return nil, nil
	}

	if l := len(fields); l != 3 && l != 4 {
		return nil, fmt.Errorf("ads.txt has fields length is incorrect.: %s", row)
	}

	// otherwise the row is valid
	var r Record
	r.ExchangeDomain = normalize(fields[0])
	r.PublisherAccountID = normalize(fields[1])
	r.AccountType = parseAccountType(normalize(fields[2]))
	// AuthorityID is optional
	if len(fields) >= 4 {
		r.AuthorityID = normalize(fields[3])
	}
	return &r, nil
}

func Parse(in io.Reader) ([]Record, error) {
	records := make([]Record, 0)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		text := scanner.Text()

		// blank line
		if len(text) == 0 {
			continue
		}

		// comment out
		if []rune(text)[0] == '#' {
			continue
		}

		// otherwise try parsing
		r, err := parseRow(scanner.Text())
		if err != nil {
			return nil, err
		}
		if r != nil {
			records = append(records, *r)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
