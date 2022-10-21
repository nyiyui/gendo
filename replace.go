package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Mappings []Mapping
}

type Mapping struct {
	Slug       string
	Comment    string
	Source     string
	SourcePath string
	TargetGlob string
}

const startKeyF = "<!-- GENDO START %s -->"
const endKeyF = "<!-- GENDO END %s -->"

func main() {
	var root, configPath string
	var config Config
	flag.StringVar(&root, "dir", "", "root directory to start replacing")
	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Parse()

	encoded, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(encoded, &config)
	if err != nil {
		log.Fatal(err)
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
		dfs := os.DirFS(root)
		matches, err := fs.Glob(dfs, mapping.TargetGlob)
		if err != nil {
			log.Fatalf("mapping %d: %s", i, err)
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
					dest.Write(source)
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
