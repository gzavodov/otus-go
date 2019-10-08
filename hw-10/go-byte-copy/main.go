package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gzavodov/otus-go/dd"
)

func main() {
	from := flag.String("from", "", "Path to source file")
	to := flag.String("to", "", "Path to destination file")
	offset := flag.Int("offset", 0, "Skips count of bytes from the start")
	limit := flag.Int("limit", 0, "Length of bytes to copy")

	flag.Parse()

	if *from == "" || *to == "" {
		fmt.Println("Usage:", os.Args[0], "-from [source file]", "-to [destination file]")
		fmt.Println("Help:", os.Args[0], "-help")
		return
	}

	copier := dd.DataCopier{}
	err := copier.Copy(*from, *to, int64(*offset), int64(*limit))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
