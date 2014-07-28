package resource

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// as this is read-only I'm not going to but a mutex on it
// I wonder how big this map will get...?
type resourceCache struct {
	resources map[string]*resourceWithTime
}

var (
	cache = resourceCache{
		resources: make(map[string]*resourceWithTime),
	}
)

type resourceWithTime struct {
	modTime time.Time
	data    []byte
}

func RegisterResource(resourceName string, resourceData []byte) {
	cache.resources[resourceName] = &resourceWithTime{time.Now(), resourceData}
}

// GetResource returns the data associated wth a resource.
// NOTE: resourceName must be a fully qualified path.
func GetResource(resourceName string) ([]byte, error) {
	if resource, ok := cache.resources[resourceName]; ok {
		if stat, err := os.Stat(resourceName); !os.IsNotExist(err) &&
			!stat.ModTime().After(resource.modTime) {
			return resource.data, nil
		}
	}

	// TODO(ttacon): hmm, this half makes me think the map needs a mutex...
	// it didn't exist in the cache, so lets load it from disk
	// we hit this if either the resource hasn't been embedded yet or if
	// the resource has been updated since the last time we used it
	data, err := ioutil.ReadFile(resourceName)
	if err != nil {
		return nil, err
	}
	cache.resources[resourceName] = &resourceWithTime{time.Now(), data}

	return data, hadToReadFromDisk
}

var hadToReadFromDisk = errors.New("had to read from disk")

var gopath = os.Getenv("GOPATH")

func GetResourceFromPkg(resourcePath, pkgPath string, isMain bool) ([]byte, error) {
	fullPath := filepath.Join(gopath, "src", pkgPath, resourcePath)
	data, err := GetResource(fullPath)
	if !isHadToReadFromDisk(err) {
		return data, err
	}

	// TODO(ttacon): be able to differentiate normal package from a main one
	pkgName := pkgPath
	if isMain {
		pkgName = "main"
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "// THIS FILE IS AUTO-GENERATED FROM %s\n", fullPath)
	fmt.Fprintf(&b, "// DO NOT EDIT.\n")
	fmt.Fprintf(&b, "package %s\n", pkgName)
	fmt.Fprintf(&b, "import \"github.com/ttacon/resource\"\n\n")
	fmt.Fprintf(&b, `
func init() {
    resource.RegisterResource(%q, []byte(%q))
}
`, fullPath, data)
	embedFileName := fmt.Sprintf("zresource_%s.go", strings.Split(
		filepath.Base(resourcePath), ".")[0])
	if dir := filepath.Dir(resourcePath); len(dir) > 0 {
		embedFileName = filepath.Join(dir, embedFileName)
	}
	embedPath := filepath.Join(gopath, "src", pkgPath, embedFileName)
	if err = ioutil.WriteFile(embedPath, b.Bytes(), 0644); err != nil {
		return nil, err
	}

	return data, nil
}

func isHadToReadFromDisk(err error) bool {
	return err == hadToReadFromDisk
}

// A Resource represents a resource on the system that is
// relative to the specified go package.
type Resource struct {
	Path       string
	RelativeTo string
	IsMain     bool
}

func (r Resource) Get() ([]byte, error) {
	return GetResourceFromPkg(r.Path, r.RelativeTo, r.IsMain)
}
