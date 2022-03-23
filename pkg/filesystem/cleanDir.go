package filesystem

import (
	"fmt"
	"path/filepath"
)

const terragruntCacheDir = ".terragrunt-cache"

func CleanDir(dir string) string {
	var output string
	fs := NewOsFs()
	if err := fs.RemoveDir(filepath.Join(dir, terragruntCacheDir)); err != nil {
		output = fmt.Sprintf("Dir: [%s] error on delete: %s", filepath.Join(dir, terragruntCacheDir), err)
	}
	return output
}
