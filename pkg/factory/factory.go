package factory

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"log"
	"os"
)

// SaveToNewFile Saves Body Content to a new file
// and prints out the parsed content to the console
func SaveToNewFile(fileName string, body *hclwrite.Body) {
	fw, osErr := os.Create(fileName)

	if osErr != nil {
		log.Fatal("Cannot create new file ", osErr)
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	for _, block := range body.Blocks() {
		log.Println(block.Type())
		for _, label := range block.Labels() {
			log.Println(label)
		}
		for key, attr := range block.Body().Attributes() {
			var tokens hclwrite.Tokens
			tokens = attr.Expr().BuildTokens(tokens)

			log.Printf("attr[%s]=%s", key, string(tokens.Bytes()))
		}
		rootBody.AppendBlock(block)
	}

	_, writeErr := newFile.WriteTo(fw)

	if writeErr != nil {
		log.Fatal("Cannot write to the new file ", writeErr)
	}
}
