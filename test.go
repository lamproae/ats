package main

import (
	"log"
	"os"
	"path/filepath"
)

func walkFunc(path string, info os.FileInfo, err error) error {
	log.Println(path, info.IsDir())
	return nil
}

func main() {
	filepath.Walk("cases", walkFunc)
}
