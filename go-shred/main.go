package main

import (
	"flag"
	"log"

	"github.com/ardelean-calin/shred"
)

func main() {
	flag.Parse()

	path := flag.Arg(0)

	err := shred.Shred(path)
	if err != nil {
		log.Fatal(err)
	}
}
