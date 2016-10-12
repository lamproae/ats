package main

import (
	"os"
	"log"
	"flag"
	"path/filepath"
	"atscase"
)

var CaseDir string = "cases"
var userCase = flag.String("c", "", "Which Case to run ? (Default All)")
var Help = flag.String("h", "help", "show help")

var CaseDB = make([]string, 0, 1000)

func WalkFunc (path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		_, file := filepath.Split(path)
		log.Println(path, file)
		_, err := os.Stat(path+"/"+file+".json")
		if err != nil {
			return nil
		}
		CaseDB = append(CaseDB, path)
	}
	return nil
}

func main() {
	flag.Parse()

	info, err := os.Stat(CaseDir+"/"+*userCase)
	if err != nil {
		log.Println("Case: ", *userCase, " is not exist")
		os.Exit(-1)
	}

	log.Println(info.Name(), info.IsDir())

	if info.IsDir() {
		filepath.Walk(CaseDir+"/"+*userCase, WalkFunc)
	}

	log.Println(CaseDB)

	for _, c := range CaseDB {
		atscase.New(c).Run()
	}
}
