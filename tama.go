package main

import (
	"flag"
	"fmt"
	"os"
)

var evalCode = flag.String("e", "", "Code for evaluation")

func usage() {
	flag.PrintDefaults()
	os.Exit(0)
}

func eval(code string) string {
	return code
}

func main() {
	flag.Parse()
	if *evalCode == "" {
		usage()
	}
	fmt.Println(eval(*evalCode))
}
