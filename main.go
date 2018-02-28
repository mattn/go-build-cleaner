package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	cis     = []*cleanerInfo{}
	dryrun  = flag.Bool("dryrun", true, "dryrun")
	verbose = flag.Bool("verbose", true, "verbose")
	cleaner = flag.String("cleaner", "*", "specify cleaners (comma separated, ? for list)")
)

type cleanerInfo struct {
	name    string
	cleaner func(bool, bool) (string, error)
	result  string
}

func register(name string, cleaner func(bool, bool) (string, error)) {
	cis = append(cis, &cleanerInfo{name: name, cleaner: cleaner})
}

func willdo(ci *cleanerInfo) bool {
	if *cleaner == "*" {
		return true
	}
	for _, name := range strings.Split(*cleaner, ",") {
		if name == ci.name {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()

	if *cleaner == "?" {
		for _, ci := range cis {
			fmt.Println(ci.name)
		}
		os.Exit(1)
	}

	for _, c := range cis {
		if !willdo(c) {
			continue
		}
		result, err := c.cleaner(*dryrun, *verbose)
		if err != nil {
			log.Printf("%s: %s\n", c.name, err)
			continue
		}
		c.result = result
	}
	for _, c := range cis {
		if c.result != "" {
			fmt.Printf("%s: %s\n", c.name, c.result)
		}
	}
	if *dryrun {
		fmt.Printf("run '%s -dryrun=false' to delete\n", os.Args[0])
	}
}
