package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	version := flag.Bool("version", false, "display version information")

	flag.Parse()

	if *version {
		fmt.Println("ased version 0.0.1")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("scripts", args[0])
	fmt.Println("files", args[1:])
}

