package filesystem

import (
	"path/filepath"
	"regexp"
)

var tfFileNameRegexp = regexp.MustCompile(`(KOMUEB|PDEKOM)-\d+\.tf$`)

func GetTerraformFileNameList(dirname string) ([]string, error) {

	fs := NewOsFs()
	fi, err := fs.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	fList := make([]string, 0)

	for _, f := range fi {
		name := f.Name()
		if tfFileNameRegexp.MatchString(name) {
			fList = append(fList, filepath.Join(dirname, name))
		}
	}

	return fList, nil
}
