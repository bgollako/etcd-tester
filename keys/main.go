package main

import (
	"etcd-tester/keys/generate"
	"etcd-tester/utils"
	"flag"
	"fmt"
)

const (
	Keys      = "keys.txt"
	DefaultN  = 1000
	OutputDir = "keys"
)

func main() {
	numKeys := flag.Int("numKeys", DefaultN, "number of unique keys to generate (default: 1000)")
	flag.Parse()

	utils.Fatal(generate.GenerateKeys(*numKeys, Keys))

	fmt.Printf("Generated %d keys in %s\n", *numKeys, Keys)
}
