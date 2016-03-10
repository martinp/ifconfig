package main

import (
	flags "github.com/jessevdk/go-flags"

	"log"
	"os"

	"github.com/martinp/ifconfigd/api"
)

func main() {
	var opts struct {
		DBPath        string `short:"f" long:"file" description:"Path to GeoIP database" value-name:"FILE" default:""`
		Listen        string `short:"l" long:"listen" description:"Listening address" value-name:"ADDR" default:":8080"`
		CORS          bool   `short:"x" long:"cors" description:"Allow requests from other domains"`
		ReverseLookup bool   `short:"r" long:"reverselookup" description:"Perform reverse hostname lookups"`
		TestIP        bool   `short:"s" long:"testip" description:"Enable IP reachability testing"`
		Template      string `short:"t" long:"template" description:"Path to template" default:"index.html"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	var a *api.API
	if opts.DBPath == "" {
		a = api.New()
	} else {
		a, err = api.NewWithGeoIP(opts.DBPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	a.CORS = opts.CORS
	a.ReverseLookup = opts.ReverseLookup
	a.Template = opts.Template
	a.TestIP = opts.TestIP

	log.Printf("Listening on %s", opts.Listen)
	if err := a.ListenAndServe(opts.Listen); err != nil {
		log.Fatal(err)
	}
}
