package sitedata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	"github.com/julienschmidt/httprouter"
)

var debug bool
var debugRootDir string

func SetConfig(d bool, drd string) {
	debug = d
	debugRootDir = drd
}

func GetStaticHandler(prefix string) httprouter.Handle {
	prefix = strings.TrimPrefix(prefix, "/")
	embedFS := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  func(path string) ([]string, error) { return nil, os.ErrNotExist },
		AssetInfo: AssetInfo,
		Prefix:    prefix,
	}
	embedHandler := http.StripPrefix("/"+prefix+"/", http.FileServer(embedFS))

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if debug {
			fs := http.Dir(filepath.Join(debugRootDir, prefix))
			handler := http.StripPrefix("/"+prefix+"/", http.FileServer(fs))
			handler.ServeHTTP(w, r)
		} else {
			embedHandler.ServeHTTP(w, r)
		}
	}
}

func AddRoutes(r *httprouter.Router, prefix string) {
	r.GET(prefix+"/*filepath", GetStaticHandler(prefix))
}

func LoadFilesInDir(dir string) (map[string]string, error) {
	dirData := map[string]string{}
	if debug {
		fullDir := filepath.Join(debugRootDir, dir)
		files, err := ioutil.ReadDir(fullDir)
		if err != nil {
			return dirData, fmt.Errorf("Error reading dir %v, %w", debugRootDir, err)
		}
		for _, file := range files {
			data, err := ioutil.ReadFile(filepath.Join(fullDir, file.Name()))
			if err != nil {
				return dirData, fmt.Errorf("Error loading %v, %w", file.Name(), err)
			}
			dirData[file.Name()] = string(data)
		}
	} else {
		files, err := AssetDir(dir)
		if err != nil {
			return dirData, fmt.Errorf("Could not load bindata dir %v, %w", dir, err)
		}
		for _, file := range files {
			fullName := path.Join("templates", file)
			data, err := Asset(fullName)
			if err != nil {
				return dirData, fmt.Errorf("Error loading bindata %v, %w", fullName, err)
			}
			dirData[file] = string(data)
		}
	}
	return dirData, nil
}
