package adstxt

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Parser is an iterative ads.txt parser
type Parser struct {
	scanner *bufio.Scanner
}

// NewParser returns a Parser
func NewParser(r io.Reader) *Parser {
	return &Parser{bufio.NewScanner(r)}
}

// Parse returns a *Record or an error
func (p *Parser) Parse() (*Record, error) {
	// scans for the first valid row or returns an error otherwise
	for p.scanner.Scan() {
		text := p.scanner.Text()

		// blank line
		if len(text) == 0 {
			continue
		}

		// comment out
		if []rune(text)[0] == '#' {
			continue
		}

		// returns when either is non-nil
		r, err := parseRow(text)
		if r != nil || err != nil {
			return r, err
		}
	}

	if err := p.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, io.EOF
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

	// if the first field contains "=", then the row is for variable declaration
	if strings.Index(fields[0], "=") != -1 {
		return nil, nil
	}

	if l := len(fields); l != 3 && l != 4 {
		return nil, fmt.Errorf("ads.txt has fields length is incorrect.: %s", row)
	}

	// otherwise the row is valid
	var r Record
	r.ExchangeDomain = strings.ToLower(normalize(fields[0]))
	r.PublisherAccountID = normalize(fields[1])
	r.AccountType = parseAccountType(normalize(fields[2]))
	// AuthorityID is optional
	if len(fields) >= 4 {
		r.AuthorityID = normalize(fields[3])
	}
	return &r, nil
}
