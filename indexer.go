package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var files = make(map[string]FileDef)
var doubleFiles = make(map[string]FileDef)

var rootPath = flag.String("rootPath", ".", "Root path for scan")
var showDoubles = flag.Bool("showDoubles", true, "Show list of double files")

type FileDef struct {
	Info *os.FileInfo
	Path string
	Hash []byte
}

func walk(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		hash := sha256.New()
		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		hash.Write(file)
		shaHash := hash.Sum(nil)
		if (*showDoubles) {
			if doubleFilePath := files[fmt.Sprintf("%x", shaHash)].Path; doubleFilePath != "" {
				fmt.Printf("%s \t -> \t %s\n", doubleFilePath, path)
				doubleFiles[fmt.Sprintf("%x", shaHash)] = FileDef{&info, path, shaHash}
			}
		} else {
			fmt.Printf("\n Indexing %s ", path)
		}
		files[fmt.Sprintf("%x", shaHash)] = FileDef{&info, path, shaHash}
	}
	return err
}

func init() {
	flag.Parse()
}
func main() {
	err := filepath.Walk(*rootPath, walk)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("\r Scan completed. Indexed %d file(s). %d duplicate(s) found", len(files), len(doubleFiles))

	os.Exit(0)
}
