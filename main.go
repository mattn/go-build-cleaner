package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	cleaners = []*cleanerInfo{}
	dryrun   = flag.Bool("dryrun", false, "dryrun")
	verbose  = flag.Bool("verbose", false, "verbose")
)

type cleanerInfo struct {
	name    string
	cleaner func(bool, bool) (string, error)
	result  string
}

func register(name string, cleaner func(bool, bool) (string, error)) {
	cleaners = append(cleaners, &cleanerInfo{name: name, cleaner: cleaner})
}

func main() {
	flag.Parse()

	for _, c := range cleaners {
		result, err := c.cleaner(*dryrun, *verbose)
		if err != nil {
			log.Fatalf("%s: %s", c.name, err)
		}
		c.result = result
	}
	for _, c := range cleaners {
		fmt.Printf("%s: %s\n", c.name, c.result)
	}
}
