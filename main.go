package main

import (
	"log"
	"os"

	"github.com/vanillaiice/gover/v3/cmd"
)

func main() {
	if err := cmd.Exec(os.Args); err != nil {
		log.Fatal(err)
	}
}
