// Code generated for package data by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/index.html
package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x55\xe1\x6a\xe4\x36\x10\xfe\xbf\x4f\x31\xa7\xfb\xd3\x42\x64\x11\x5a\xca\xd1\xda\x0b\xd7\x5c\x29\xa5\x81\x42\xb9\x17\x18\x4b\xb3\xb6\x36\xb2\xe4\x93\xc6\x9b\x5d\x42\xde\xbd\xd8\x6b\xe7\x6c\xef\x36\x29\x2d\x34\x10\xd6\x1a\x8f\xbe\xf9\x66\xe6\x9b\x71\xfe\xce\x04\xcd\xa7\x96\xa0\xe6\xc6\x6d\x37\x9b\xbc\x26\x34\xdb\x0d\x00\x40\xde\x10\x23\xe8\x1a\x63\x22\x2e\x44\xc7\x3b\xf9\x41\xcc\x5f\x79\x6c\xa8\x10\x07\x4b\x8f\x6d\x88\x2c\x40\x07\xcf\xe4\xb9\x10\x8f\xd6\x70\x5d\x18\x3a\x58\x4d\x72\x38\xdc\x80\xf5\x96\x2d\x3a\x99\x34\x3a\x2a\x6e\x27\x20\x67\xfd\x03\x44\x72\x85\x48\x7c\x72\x94\x6a\x22\x16\x50\x47\xda\x15\xa2\x66\x6e\xd3\x8f\x4a\x35\x78\xd4\xc6\x67\x65\x08\x9c\x38\x62\xdb\x1f\x74\x68\xd4\x8b\x41\x7d\x97\x7d\x9f\xdd\x2a\x9d\xd2\x57\x5b\xd6\x58\x9f\xe9\x94\xa6\x40\x49\x47\xdb\x32\xa4\xa8\xbf\x02\xe3\x1e\x8f\x59\x15\x42\xe5\x08\x5b\x9b\x06\xd0\xde\xa6\x9c\x2d\x93\xda\x7f\xe9\x28\x9e\x46\xec\xf3\x61\x00\xdd\x27\xb1\xcd\xd5\x19\xef\x15\xf0\x7f\xca\x7a\xbf\x26\x7d\x05\xff\x9d\x94\x43\xc5\x7b\x6c\x49\x5f\x3a\x7b\x28\x44\xa4\x5d\xa4\x54\xcf\xca\xfe\xe1\xa7\x2e\xba\x62\x8a\x5f\x22\x6b\x3c\x50\xb6\xa7\xd8\xa0\xcf\xac\xdf\x05\x25\xe5\x08\xc8\x96\x1d\x6d\x7f\xc5\x88\x15\xc1\xa7\x10\x62\xae\xce\xa6\x4d\xae\xce\x0a\xd8\xe4\x65\x30\xa7\xd1\xdd\xe3\x01\xb4\xc3\x94\x0a\xe1\xf1\x50\x62\x84\xf3\x8f\x34\xb4\xc3\xce\xf1\x74\xdc\xd9\x23\x19\xc9\xa1\x1d\xab\x3e\x5c\x36\xf6\xe5\x72\x4f\x15\xad\xa7\x28\x77\xae\xb3\x66\xe6\x75\x4e\x12\x7e\x8e\xe8\x0d\xf4\xff\x1c\xaa\xca\x11\x54\xc4\x50\xc5\xd0\xb5\x64\x60\x17\x22\x94\xc4\x4c\x11\x9a\x50\x5a\x47\x60\x6c\x6a\x1d\x9e\x60\xca\xeb\x32\xe2\x48\xac\xcf\x89\xe2\x22\x1e\x40\x5e\x76\xcc\xc1\x43\xaf\xff\x42\x9c\x0f\x62\x75\x71\xa4\xa1\x83\x73\xd8\x26\x32\x02\x0c\x32\x8e\xe6\x3e\xa1\xb3\x7d\x32\x63\xac\xfa\x59\x79\x5f\x26\x49\x47\x6c\x5a\x47\x72\x04\x9a\x3c\xe5\xad\x00\x8c\x16\x25\x1d\x5b\xf4\x86\x4c\x21\x76\xe8\x12\xad\xb8\xf5\xa2\x6a\xd1\x4f\x6c\x52\x94\xc1\xbb\x93\xd8\x7e\x3e\xf3\xf1\x78\xb0\x15\xb2\x0d\x3e\x57\xbd\xdf\xab\x97\xad\x0e\x5e\x96\x18\x07\x59\xfd\x8f\xce\xb9\x3a\xd7\x74\x65\xc5\x55\x89\xcb\xbe\xe7\xd3\xcc\xbf\x17\x4b\x55\xe2\xa2\xb1\xca\xd8\xc3\x64\x58\x4b\xe7\x2e\x38\x47\x9a\x81\xeb\xa1\x3c\xd0\x6f\x96\x74\xd3\x8b\xa6\x49\x37\x83\xa4\x02\xd7\x14\xa7\x71\x19\xd4\x34\xb4\xd1\xfa\xea\x15\x01\x4d\x7d\x83\x55\x1f\x05\x58\x53\x88\x57\xfb\xbc\xca\xbb\x73\xb3\xc4\x27\x38\x8f\x87\xcb\xce\x3b\x3b\x79\xa2\x66\x7b\xb8\xd4\x46\xff\xf7\xf4\x04\x76\x07\x99\x0b\x95\xf5\xf0\xfc\x7c\xc5\x23\xc7\xb1\xa8\x0a\x3b\xae\xd5\xe0\x29\xb6\xf7\xc3\x85\xeb\xea\xfa\x46\x77\x31\x92\xe7\x6f\xc7\x76\xae\xea\x3f\x8b\x4d\x2e\xd1\xd5\xa0\x57\xa2\x86\x8e\x87\xb0\xa1\x63\x78\x7a\xca\xba\x44\xb1\xff\x74\x3c\x3f\xff\x67\x12\xde\x5c\xe1\x90\x2b\x67\x2f\x4a\xaa\x3a\xb7\xb4\x5d\xc8\x6a\x10\x91\xca\x56\x6d\x5c\x28\x63\xe1\xb8\xda\x65\x33\xc7\x5c\x79\x5c\xcb\x74\x2e\xa9\x7d\xd7\x94\x81\x63\xbf\x79\xe8\xc8\x52\x93\xe7\xc5\x6e\xca\xeb\xdb\xe5\x10\xd4\xb7\xb3\x97\xed\xf6\xe3\x50\x46\xb6\x43\x09\x55\xbb\x78\xf7\xb9\x26\x30\x21\x44\x48\x8c\x4c\x60\x53\xef\x3a\x3c\xcf\x7c\xc7\x39\xda\xac\x89\xbd\xa4\xf4\x37\xeb\x3b\x86\xc7\xb5\xa8\x97\x93\x22\x7f\xb8\x22\xd5\x99\x20\x5a\xab\xb4\xb3\xfa\x41\x40\x0c\xee\x72\xe7\x96\xec\xa1\x64\x2f\x53\xa7\x35\xa5\x34\x3c\x97\x2e\xe8\x87\xab\x03\x00\x70\x77\xff\xdb\xdd\xef\x17\xf2\x58\x6c\x89\x7f\x43\xf3\x2d\x76\xa4\x83\x37\x18\x4f\x6f\xf2\xfb\xf3\x97\xfb\x3f\x3e\x7e\x7a\x8b\xe0\xec\x38\x3e\xe6\x6a\xf8\xf6\xfe\x15\x00\x00\xff\xff\x69\xf6\xd4\xe7\x9d\x09\x00\x00")

func templatesIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesIndexHtml,
		"templates/index.html",
	)
}

func templatesIndexHtml() (*asset, error) {
	bytes, err := templatesIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/index.html", size: 2461, mode: os.FileMode(436), modTime: time.Unix(1588557920, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/index.html": templatesIndexHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"index.html": &bintree{templatesIndexHtml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
