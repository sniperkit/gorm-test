// CSV to JSON converter.
// NOTE: JSON keys must not contain '-' because it's used as a slice indicator.

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iharsuvorau/json2csv2json/csvjson"
)

func main() {
	fileOrDir := flag.String("i", ".", "input file or directory")
	outputDir := flag.String("o", "./", "output directory to save files to")
	delimeter := flag.String("d", "\t", "delimeter for the input file: \t or ,")
	flag.Parse()

	f, err := os.Open(*fileOrDir)
	if err != nil {
		log.Fatal(err)
	}

	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// define files' extenstions from the delimeter
	var outputExtenstion, outputFilename string
	var delimeterRune rune

	outputExtenstion = ".json"

	switch *delimeter {
	case "\t":
		delimeterRune = '\t'
	case ",":
		delimeterRune = ','
	}

	if stat.IsDir() {
		// walk the directory recursively and look for input files
		if err = filepath.Walk(*fileOrDir, csvjson.WalkWokring(*outputDir, outputExtenstion, delimeterRune)); err != nil {
			log.Fatal(err)
		}
	} else {
		// replace the extension
		outputFilename = strings.Replace(filepath.Base(*fileOrDir), filepath.Ext(*fileOrDir), "", 1) + outputExtenstion

		// do the work
		if err = csvjson.Work(*fileOrDir, filepath.Join(*outputDir, outputFilename), delimeterRune); err != nil {
			log.Fatal(err)
		}
	}
}
