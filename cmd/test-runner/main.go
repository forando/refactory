package main

import (
	"flag"
	"github.com/forando/refactory/pkg/factory"
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/parser"
	"github.com/forando/refactory/pkg/schema"
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
	dirFlag := flag.String("dir", "/Users/andrii.logoshko/Projects/aws-accounts/aws-prod-org", "path to a dir to scan")

	flag.Parse()

	if len(*dirFlag) == 0 {
		log.Fatal("requires path to a dir. please use the -dir flag to set it")
	}

	fList, err := filesystem.GetTerraformFileNameList(*dirFlag)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fList {
		log.Printf("File %s:", file)
		body := parser.ParseFile(file)
		permissionSetNames := make(map[string]string)
		var permissionSetModules schema.PermissionSetModules
		for _, block := range body.Blocks() {
			blockMetaData, err := parser.ParseBlockType(block)
			if err != nil {
				log.Fatalf("File %s: %s", file, err)
			}
			if blockMetaData.BlockType == schema.PermissionSetModuleType {
				module, err := parser.ParsePermissionSetModule(block.Body())
				if err != nil {
					log.Fatal(err)
				}
				permissionSetModules = append(permissionSetModules, module)
				permissionSetNames[blockMetaData.BlockName] = module.PermissionSetName
			}
		}
		log.Println(permissionSetNames)

		var accountModules schema.AccountModules
		for _, block := range body.Blocks() {
			blockMetaData, err := parser.ParseBlockType(block)
			if err != nil {
				log.Fatalf("File %s: %s", file, err)
			}
			if blockMetaData.BlockType == schema.AccountModuleType {
				module, err := parser.ParseAccountModule(block.Body(), permissionSetNames)
				if err != nil {
					log.Fatal(err)
				}
				accountModules = append(accountModules, module)
				log.Println(module)
			}
		}
		//factory.Bootstrap(&accountModules)
	}
}
