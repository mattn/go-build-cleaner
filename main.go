package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	cleaners = []*cleanerInfo{}
	dryrun   = flag.Bool("dryrun", true, "dryrun")
	verbose  = flag.Bool("verbose", true, "verbose")
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
			log.Printf("%s: %s\n", c.name, err)
			continue
		}
		c.result = result
	}
	for _, c := range cleaners {
		fmt.Printf("%s: %s\n", c.name, c.result)
	}
	if *dryrun {
		fmt.Printf("run '%s -dryrun=false' to delete\n", os.Args[0])
	}
}
