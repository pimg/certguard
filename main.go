package main

import (
	"log"

	"github.com/pimg/crl-inspector/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
