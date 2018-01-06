package main

import (
	"flag"
	"os"

	logger "github.com/mgamsjager/indexer/log"
	"github.com/mgamsjager/indexer/search"
)

func main() {
	var rootPath = flag.String("root-path", ".", "Root path for scan")
	var deleteDuplicates = flag.Bool("delete-duplicates", false, "Delete found duplicates from file system")
	var sampleSize = flag.Int64("sample-size", 100*1024, "Number of bytes to hash.")

	flag.Parse()
	s := search.NewSearcher(search.Config{
		Logger:     logger.New(),
		Path:       *rootPath,
		SampleSize: *sampleSize,
		Delete:     *deleteDuplicates,
	})

	if err := s.FindDuplicates(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
