# Indexer

Indexer is a simple CLI tool which search for duplicate files by
comparing its' SHA256 hash.

# Usage
Run `$ indexer` to give you all duplicates in the current and descending directories

## Flags
`-root-path` Specify the path to start from

`-delete-duplicates` Will automatically delete found duplicates!! 

`-sample-size` Number of bytes to calculate the hash from. Default is 4000. 