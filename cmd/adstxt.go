// adstxt cralwer implementation
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/suzuken/go-adstxt"
)

func main() {
	var (
		rawurl = flag.String("url", "", "URL of ads.txt to crawl")
	)
	flag.Parse()
	ads, err := adstxt.Get(*rawurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read ads.txt failed: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%#v", ads)
}
