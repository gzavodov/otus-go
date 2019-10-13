package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gzavodov/otus-go/envdir"
)

func main() {

	path := ""
	executable := ""
	if len(os.Args) >= 3 {
		path = os.Args[1]
		executable = os.Args[2]
	}

	if path == "" || executable == "" {
		fmt.Println("Usage:", os.Args[0], "<environment directory>", "<executable>")
		return
	}

	envmodifier := envdir.EnvModifier{}
	_, err := envmodifier.Run(path, executable)
	if err != nil {
		log.Fatal(err)
	}
}
