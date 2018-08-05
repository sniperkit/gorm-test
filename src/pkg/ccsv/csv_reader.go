package ccsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

type CsvReader struct {
	File  string
	Keys  []string
	Items []map[string]string
	file  *os.File
}

func (c *CsvReader) parseCSVData(parseItems chan []string) {
	csvfile, err := os.Open(c.File)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1
	rowCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if rowCount == 0 {
			c.Keys = record
		}
		rowCount++
		// Stop at EOF.
		go func(pi []string) {
			parseItems <- pi
		}(record)
	}
}

func (c *CsvReader) moldObject(parseItems <-chan []string, wg *sync.WaitGroup, lineItems chan<- map[string]string) {
	defer wg.Done()
	for pi := range parseItems {
		var l map[string]string
		for i, key := range c.Keys {
			fmt.Println(l[key])
			fmt.Println(pi[i])
			fmt.Println(i)
			l[key] = pi[i]
		}
		lineItems <- l
	}
}

func (c *CsvReader) Run() {
	parseItems := make(chan []string)
	lineItems := make(chan map[string]string)

	go c.parseCSVData(parseItems)

	wg := new(sync.WaitGroup)
	for i := 0; i <= 3; i++ {
		wg.Add(1)
		go c.moldObject(parseItems, wg, lineItems)
	}

	go func() {
		wg.Wait()
		close(lineItems)
	}()

	for li := range lineItems {
		c.Items = append(c.Items, li)
	}
}
