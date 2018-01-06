# Indexer

Indexer is a simple CLI tool which search for duplicate files by
comparing its' SHA256 hash.

# Install

`$ go install github.com/mgamsjager/indexer/cmd/indexer`

# Usage
Run `$ indexer` to give you all duplicates in the current and descending directories

## Flags
`-root-path` Specify the path to start from

`-delete-duplicates` Will automatically delete found duplicates!! 

`-sample-size` Number of bytes to calculate the hash from. Default is 100KB. 
