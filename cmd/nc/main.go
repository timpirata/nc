package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	serveHTTP  = flag.Bool("serve", false, "Run as webserver / cloud function")
	AppVersion string
)

func main() {
	if os.Getenv("NC_SERVE_HTTP") != "" {
		*serveHTTP = true
	}
	flag.Parse()

	httpdConfig := httpdConfig{
		address: thisOrThat(os.Getenv("NC_TCP_PORT"), ":2001"),
		enabled: *serveHTTP,
	}
	fmt.Printf("NC version %s, inspired by Peter Norton et al, arising ..\n", AppVersion)
	RunWebserver(httpdConfig)
	// ProducePDF
	// ProduceTEX
	fmt.Println("My job is done here - pleased to meet ya!")
}
