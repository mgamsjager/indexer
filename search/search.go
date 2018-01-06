package search

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/mgamsjager/indexer/log"
)

var files = make(map[string]fileDef)

type searcher struct {
	Config
	Counter *counter
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

type counter struct {
	Files int64
	Mutex sync.RWMutex
}

func NewSeacher(c Config) *searcher {
	return &searcher{c, &counter{0, sync.RWMutex{}}}
}

func (s *searcher) FindDuplicates() error {
	s.Logger.Infof("Starting to scan from %s \n ", s.Config.Path)
	if err := filepath.Walk(s.Config.Path, s.walk); err != nil {
		return err
	}
	s.Logger.Infof("\nDone! \t Files scanned: %d\n", s.Counter.Files)
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

	if file == nil {
		return nil
	}

	hash.Write(file)
	shaHash := hash.Sum(nil)

	if doubleFilePath := checkForDuplicate(shaHash); doubleFilePath != "" {
		s.Logger.Infof("\n%s \t -> \t %s\n", doubleFilePath, path)
		if s.Config.Delete {
			delete(path)
		}
	}
	registerFile(path, shaHash)
	s.counter()
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
			s.Logger.Error(err)
			return nil, err
		}
		return b, nil
	} else {
		buf, err := ioutil.ReadAll(f)
		return buf, err
	}
	return nil, err
}

func (s *searcher) counter() {
	s.Counter.Mutex.Lock()
	s.Counter.Files += 1
	s.Logger.Infof("\r Scanned %d files", s.Counter.Files)
	s.Counter.Mutex.Unlock()
}
