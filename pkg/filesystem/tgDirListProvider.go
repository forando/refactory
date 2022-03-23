package filesystem

import (
	"regexp"
)

var tgDirNameRegexp = regexp.MustCompile(`(KOMUEB|PDEKOM)-\d+$`)
var tgPermissionSetsDirRegexp = regexp.MustCompile(`PermissionSets$`)

func GetTerragruntModuleNameList(dirname string) (*[]string, error) {

	fs := NewOsFs()
	dirs, err := fs.ListDirs(dirname)
	if err != nil {
		return nil, err
	}

	productDirs := make([]string, 0)

	for _, dir := range dirs {
		if tgDirNameRegexp.MatchString(dir) {
			productDirs = append(productDirs, dir)
		}
	}

	imports := make([]string, 0)
	for _, productDir := range productDirs {
		productModules, err := fs.ListDirs(productDir)
		if err != nil {
			return nil, err
		}
		for _, productModule := range productModules {
			if tgPermissionSetsDirRegexp.MatchString(productModule) {
				permissionSets, err := fs.ListDirs(productModule)
				if err != nil {
					return nil, err
				}
				for _, pSet := range permissionSets {
					imports = append(imports, pSet)
				}
			} else {
				imports = append(imports, productModule)
			}
		}
	}

	return &imports, nil
}
