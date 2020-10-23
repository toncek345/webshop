package main

import (
	"log"
	"os"

	"github.com/toncek345/webshop/cmd"
)

func main() {
	if err := cmd.RegisterCmds().Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
