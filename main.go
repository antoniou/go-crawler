package main

import (
	"log"
	"os"

	"github.com/antoniou/go-crawler/client"
)

func main() {
	err := client.New().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
