package factory

import (
	"encoding/csv"
	"github.com/forando/refactory/pkg/schema"
	"log"
	"os"
	"path"
)

func BootstrapImports(dir string, module string, importResources *map[string][]schema.TfImport) {

	var imports []schema.TfImport
	found := false
	if imports, found = (*importResources)[module]; !found || len(imports) == 0 {
		log.Fatalf("Cannot find importResource for module %s", module)
	}

	f, err := os.Create(path.Join(dir, "imports.csv"))
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(f)
	w.Comma = ';'
	defer w.Flush()

	for _, imp := range imports {
		record := []string{module, imp.Address, imp.Id}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
