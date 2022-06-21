package filesystem

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FS represents a minimal filesystem implementation
// See io/fs.FS in http://golang.org/s/draft-iofs-design
type FS interface {
	Open(name string) (File, error)
	RemoveDir(name string) error
	RemoveFile(name string) error
	ReadFile(name string) ([]byte, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	MakeDirs(name string)
	ListDirs(string) ([]string, error)
	Exists(name string) (bool, error)
}

// File represents an open file in FS
// See io/fs.File in http://golang.org/s/draft-iofs-design
type File interface {
	Stat() (os.FileInfo, error)
	Read([]byte) (int, error)
	Close() error
}

type osFs struct{}

func (fs *osFs) Open(name string) (File, error) {
	return os.Open(name)
}

func (fs *osFs) RemoveDir(name string) error {
	return os.RemoveAll(name)
}

func (fs *osFs) RemoveFile(name string) error {
	return os.Remove(name)
}

func (fs *osFs) ReadFile(name string) ([]byte, error) {
	//TODO: potential cause for `too many open files` error
	return ioutil.ReadFile(name)
}

func (fs *osFs) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (fs *osFs) ListDirs(dir string) (dirs []string, err error) {
	infos, err := fs.ReadDir(dir)
	if err != nil {
		return
	}

	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		path := filepath.Join(dir, info.Name())
		dirs = append(dirs, path)
	}
	return
}

func NewOsFs() FS {
	return &osFs{}
}

func DirFiles(fs FS, dir string) (primary []string, err error) {
	infos, err := fs.ReadDir(dir)
	if err != nil {
		err = errors.Errorf("Directory %s does not exist or cannot be read.", dir)
		return
	}

	var override []string
	for _, info := range infos {
		if info.IsDir() {
			continue
		}

		name := info.Name()
		ext := fileExt(name)
		if ext == "" || isIgnoredFile(name) {
			continue
		}

		baseName := name[:len(name)-len(ext)] // strip extension
		isOverride := baseName == "override" || strings.HasSuffix(baseName, "_override")

		fullPath := filepath.Join(dir, name)
		if isOverride {
			override = append(override, fullPath)
		} else {
			primary = append(primary, fullPath)
		}
	}

	// We are assuming that any _override files will be logically named,
	// and processing the files in alphabetical order. Primaries first, then overrides.
	primary = append(primary, override...)

	return
}

func (fs *osFs) MakeDirs(name string) {
	err := os.MkdirAll(name, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func (fs *osFs) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func fileExt(path string) string {
	if strings.HasSuffix(path, ".tf") {
		return ".tf"
	} else if strings.HasSuffix(path, ".tf.json") {
		return ".tf.json"
	} else {
		return ""
	}
}

func isIgnoredFile(name string) bool {
	return strings.HasPrefix(name, ".") || // Unix-like hidden files
		strings.HasSuffix(name, "~") || // vim
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#") // emacs
}

func PrintPWD() {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:")
	fmt.Println(dir)
}
