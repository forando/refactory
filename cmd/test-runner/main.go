package main

import (
	"flag"
	"github.com/forando/refactory/pkg/factory"
	"github.com/forando/refactory/pkg/parser"
	"log"
)

func main() {
	fileFlag := flag.String("file", "./inputs/simple-input.tf", "path to a file to parse")

	flag.Parse()

	if len(*fileFlag) == 0 {
		log.Fatal("requires path to a file. please use the -file flag to set it")
	}

	body := parser.ParseFile(*fileFlag)

	factory.SaveToNewFile("./output.tf", body)
}
