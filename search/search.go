package search

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mgamsjager/indexer/log"
)

var files = make(map[string]fileDef)

type searcher struct {
	Config
}

type Config struct {
	Logger     log.Logger
	Path       string
	Delete     bool
	SampleSize int64
}

type fileDef struct {
	Path string
	Hash []byte
}

func NewSeacher(c Config) *searcher {
	return &searcher{c}
}

func (s *searcher) FindDuplicates() error {
	s.Logger.Info("Starting to scan from " + s.Config.Path)
	if err := filepath.Walk(s.Config.Path, s.walk); err != nil {
		return err
	}
	s.Logger.Info("Done")
	return nil
}

func (s *searcher) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		s.Logger.Error("Read error", err)
		return err
	}

	if info.IsDir() {
		return nil
	}

	hash := sha256.New()
	file, err := s.readFile(path)
	if err != nil {
		s.Logger.Fatal(err)
	}

	hash.Write(file)
	shaHash := hash.Sum(nil)

	if doubleFilePath := checkForDuplicate(shaHash); doubleFilePath != "" {
		s.Logger.Infof("\n%s \t -> \t %s\n", doubleFilePath, path)
		if s.Config.Delete {
			delete(path)
		}
	} else {
		registerFile(path, shaHash)
	}
	return nil
}

func checkForDuplicate(shaHash []byte) string {
	var hexHash = fmt.Sprintf("%x", shaHash)
	return files[hexHash].Path
}

func registerFile(path string, shaHash []byte) {
	var hexHash = fmt.Sprintf("%x", shaHash)
	files[hexHash] = fileDef{path, shaHash}
}

func delete(path string) {
	fmt.Println("Deleting", path)
	os.Remove(path)
}

func (s *searcher) readFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if stat, err := f.Stat(); stat.Size() >= s.Config.SampleSize+1 {
		if err != nil {
			return nil, err
		}

		sr := io.NewSectionReader(f, 0, s.Config.SampleSize)

		b := make([]byte, s.Config.SampleSize)
		if _, err := sr.Read(b); err != nil {
			s.Logger.Fatal(err)
		}
		return b, nil
	} else {
		buf, err := ioutil.ReadAll(f)
		return buf, err
	}
	return nil, err
}
