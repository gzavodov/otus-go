package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gzavodov/otus-go/envdir"
)

func main() {

	var path string
	var executable string

	if len(os.Args) >= 3 {
		path = os.Args[1]
		executable = os.Args[2]
	}

	if path == "" || executable == "" {
		fmt.Println("Usage:", os.Args[0], "<environment variables directory>", "<executable>")
		return
	}

	envmodifier := envdir.EnvModifier{}
	output, err := envmodifier.Run(path, executable)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", output)
}
