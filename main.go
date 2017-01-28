package main

import (
	"os"

	"github.com/antoniou/go-crawler/client"
)

func main() {
	client.New().Run(os.Args)
}
