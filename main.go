package main

import (
	"log"

	"github.com/pimg/certguard/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
