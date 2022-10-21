package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Mappings []Mapping
}

type Mapping struct {
	Slug       string
	Comment    string
	Data       map[string]string
	Source     string
	SourcePath string
	TargetGlob string
}

const startKeyF = "<!-- GENDO START %s -->"
const endKeyF = "<!-- GENDO END %s -->"

type TmplData struct {
	Mapping Mapping
	Path    string
}

func main() {
	var root, configRaw string
	var config Config
	flag.StringVar(&root, "dir", "", "root directory to start replacing")
	flag.StringVar(&configRaw, "config", "", "raw config")
	flag.Parse()

	err := json.Unmarshal([]byte(configRaw), &config)
	if err != nil {
		log.Printf("config raw:\n%s", configRaw)
		log.Fatalf("json: %s", err)
	}

	for i, mapping := range config.Mappings {
		source := []byte(mapping.Source)
		if path := mapping.SourcePath; path != "" {
			source, err = os.ReadFile(path)
			if err != nil {
				log.Fatalf("mapping %d: %s", i, err)
			}
		}

		startKey := fmt.Sprintf(startKeyF, mapping.Slug)
		endKey := fmt.Sprintf(endKeyF, mapping.Slug)
		matches, err := filepath.Glob(filepath.Join(root, mapping.TargetGlob))
		if err != nil {
			log.Fatalf("mapping %d: %s", i, err)
		}
		tmpl, err := template.New("").Parse(string(source))
		if err != nil {
			log.Fatalf("mapping %d: tmpl: %s", i, err)
		}
		for _, path := range matches {
			dest := new(bytes.Buffer)
			src, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("mapping %d match %s: %s", i, path, err)
			}
			src2 := src
			for {
				{
					startIndex := bytes.Index(src2, []byte(startKey))
					if startIndex == -1 {
						dest.Write(src2)
						break
					}
					dest.Write(src2[:startIndex])
					src2 = src2[startIndex:]
				}
				{
					endIndex := bytes.Index(src2, []byte(endKey))
					if endIndex == -1 {
						log.Fatalf("mapping %d match %s: end key not found", i, path)
					}
					io.WriteString(dest, startKey)
					err = tmpl.Execute(dest, TmplData{
						Mapping: mapping,
						Path:    path,
					})
					if err != nil {
						log.Fatalf("mapping %d match %s: tmpl: %s", i, path, err)
					}
					src2 = src2[endIndex:]
				}
			}
			err = os.WriteFile(path, dest.Bytes(), 0)
			if err != nil {
				log.Fatalf("mapping %d match %s: %s", i, path, err)
			}
			log.Printf("mapping %d match %s", i, path)
		}
	}
}
