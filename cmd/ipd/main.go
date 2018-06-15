package main

import (
	"log"

	flags "github.com/jessevdk/go-flags"

	"os"

	"github.com/mpolden/ipd/http"
	"github.com/mpolden/ipd/iputil"
	"github.com/mpolden/ipd/iputil/database"
)

func main() {
	var opts struct {
		CountryDBPath string `short:"f" long:"country-db" description:"Path to GeoIP country database" value-name:"FILE" default:""`
		CityDBPath    string `short:"c" long:"city-db" description:"Path to GeoIP city database" value-name:"FILE" default:""`
		ASNDBPath     string `short:"a" long:"asn-db" description:"Path to GeoIP ASN database" value-name:"FILE" default:""`
		Listen        string `short:"l" long:"listen" description:"Listening address" value-name:"ADDR" default:":8080"`
		ReverseLookup bool   `short:"r" long:"reverse-lookup" description:"Perform reverse hostname lookups"`
		PortLookup    bool   `short:"p" long:"port-lookup" description:"Enable port lookup"`
		Template      string `short:"t" long:"template" description:"Path to template" default:"index.html" value-name:"FILE"`
		IPHeader      string `short:"H" long:"trusted-header" description:"Header to trust for remote IP, if present (e.g. X-Real-IP)" value-name:"NAME"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	log := log.New(os.Stderr, "ipd: ", 0)
	db, err := database.New(opts.CountryDBPath, opts.CityDBPath, opts.ASNDBPath)
	if err != nil {
		log.Fatal(err)
	}

	server := http.New(db)
	server.Template = opts.Template
	server.IPHeader = opts.IPHeader
	if opts.ReverseLookup {
		log.Println("Enabling reverse lookup")
		server.LookupAddr = iputil.LookupAddr
	}
	if opts.PortLookup {
		log.Println("Enabling port lookup")
		server.LookupPort = iputil.LookupPort
	}
	if opts.IPHeader != "" {
		log.Printf("Trusting header %s to contain correct remote IP", opts.IPHeader)
	}

	log.Printf("Listening on http://%s", opts.Listen)
	if err := server.ListenAndServe(opts.Listen); err != nil {
		log.Fatal(err)
	}
}
