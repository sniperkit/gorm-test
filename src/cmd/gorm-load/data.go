package main

import (
	// "fmt"
	// "log"
	// "runtime"

	// internal
	"github.com/sniperkit/gorm-test/src/pkg/ccsv"
)

var (
	cWriterCSV *ccsv.CsvWriter
	cReaderCSV *ccsv.CsvReader
)

func newCsvWriter(fp string) (*ccsv.CsvWriter, error) {
	return ccsv.NewCsvWriter(fp)
}

func newCsvReader(fp string) (*ccsv.CsvWriter, error) {
	return ccsv.NewCsvWriter(fp)
}

/*
func loadFileDB(filePath string) {
    wrap := "\n"
    if runtime.GOOS == "windows" {
        wrap = "\r\n"
    }
    mysql.RegisterLocalFile(filePath)
    loadSql := "load data LOW_PRIORITY local infile '%s' into table url FIELDS TERMINATED BY ','  LINES TERMINATED BY '%s' (Url,`Name`)  set UrlGroupId=1"
    loadSql = fmt.Sprintf(loadSql, filePath, wrap)
    _, err := dbs.Exec(loadSql)
    if err != nil {
        log.Fatalf("sql is warn, sql: %s err: %s \n", loadSql, err.Error())
        sqlClose()
        sqlConntion()
        return
    }
    log.Printf("db run success! filePath: %s \n", filePath)
}
*/
