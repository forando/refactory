package main

import (
	"flag"
	"github.com/forando/refactory/pkg/factory"
	"github.com/forando/refactory/pkg/parser"
	"log"
)

func main1() {
	fileFlag := flag.String("file", "./inputs/simple-input.tf", "path to a file to parse")

	flag.Parse()

	if len(*fileFlag) == 0 {
		log.Fatal("requires path to a file. please use the -file flag to set it")
	}

	body := parser.ParseFile(*fileFlag)

	factory.SaveToNewFile("./output.tf", body)
	//factory.BootstrapAccountTerragrunt("./terragrunt_1.hcl")
}

func main() {
	fileFlag := flag.String("file", "./inputs/simple-input.tf", "path to a file to parse")

	flag.Parse()

	if len(*fileFlag) == 0 {
		log.Fatal("requires path to a file. please use the -file flag to set it")
	}

	body := parser.ParseFile(*fileFlag)

	moduleBody := body.Blocks()[0].Body()

	module, err := parser.ParseAccountModule(moduleBody)
	if err != nil {
		log.Fatal(err)
	}

	/*modules := make(schema.AccountModules, 1)
	modules[0] = module

	factory.Bootstrap(&modules)*/
	log.Println(module)
}
