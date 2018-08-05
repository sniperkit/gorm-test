// JSON to CSV converter.
// NOTE: JSON keys must not contain '-' because it's used as a slice indicator.

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iharsuvorau/json2csv2json/jsoncsv"
)

func main() {
	fileOrDir := flag.String("i", ".", "input file or directory")
	outputDir := flag.String("o", "./", "output directory to save files to")
	delimeter := flag.String("d", "\t", "delimeter for the output file: \t or ,")
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

	switch *delimeter {
	case "\t":
		outputExtenstion = ".tsv"
		delimeterRune = '\t'
	case ",":
		outputExtenstion = ".csv"
		delimeterRune = ','
	}

	if stat.IsDir() {
		// walk the directory recursively and look for .json files
		if err = filepath.Walk(*fileOrDir, jsoncsv.WalkWokring(*outputDir, outputExtenstion, delimeterRune)); err != nil {
			log.Fatal(err)
		}
	} else {
		// replace the extension
		outputFilename = strings.Replace(filepath.Base(*fileOrDir), filepath.Ext(*fileOrDir), "", 1) + outputExtenstion

		// do the work
		if err = jsoncsv.Work(*fileOrDir, filepath.Join(*outputDir, outputFilename), delimeterRune); err != nil {
			log.Fatal(err)
		}
	}
}
