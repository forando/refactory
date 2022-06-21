package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

type TempFile struct {
	Name string
}

func (f *TempFile) Path() string {
	return filepath.Join(os.TempDir(), f.Name)
}

func (f *TempFile) Remove() {
	if err := os.Remove(f.Path()); err != nil {
		fmt.Printf("Could not delete temp file\nError: %s", err.Error())
	}
}
