// Copyright 2018 John Scherff
//
// Licensed under the Apache License, version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csvio

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

// NewReadWriter returns a new ReadWriter object.
func NewReadWriter() ReadWriter {
	return new(records)
}

// ReadWriter is an interface that exposes methods for reading
// JSON input and writing CSV output.
type ReadWriter interface {
	Read(io.Reader) error
	ReadUrl(string) error
	ReadFile(string) error
	ReadRecord([]string)
	Write(io.Writer) error
	WriteFile(string) error
}

// records contains a slice of string slices that represent
// one or more records of data fields collected from an array
// of JSON objects. The first record of the records contains
// the field names. The remaining records contain the values.
type records [][]string

// fields returns a sorted list of property names in a thing.
func keys(thing map[string]interface{}) (keys []string) {

	for key := range thing {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

// Read unmarshals a records of JSON objects from an []byte.
func (this *records) Read(r io.Reader) error {

	var things []map[string]interface{}

	if err := json.NewDecoder(r).Decode(&things); err != nil {
		return err
	} else if len(things) == 0 {
		return fmt.Errorf(`empty records`)
	}

	var fields []string

	for index, thing := range things {

		if index == 0 {
			fields = keys(thing)
			this.ReadRecord(fields)
		}

		var record []string

		for _, field := range fields {

			var value string

			if thing[field] != nil {
				value = fmt.Sprintf(`%v`, thing[field])
			}

			record = append(record, value)
		}

		this.ReadRecord(record)
	}

	return nil
}

// ReadUrl unmarshals a records of JSON objects from a URL.
func (this *records) ReadUrl(url string) error {

	if resp, err := http.Get(url); err != nil {
		return err
	} else {
		defer resp.Body.Close()
		return this.Read(resp.Body)
	}
}

// ReadFile unmarshals a records of JSON objects from JSON file.
func (this *records) ReadFile(fn string) error {

	if fh, err := os.Open(fn); err != nil {
		return err
	} else {
		defer fh.Close()
		return this.Read(fh)
	}
}

// ReadRecord reads a slice of strings as a record.
func (this *records) ReadRecord(record []string) {
	*this = append(*this, record)
}

// Write writes the accumulated records to an io.Writer in CSV format.
func (this *records) Write(w io.Writer) error {

	if err := csv.NewWriter(w).WriteAll(*this); err != nil {
		return err
	}

	return nil
}

// Write writes the accumulated records to a file in CSV format.
func (this *records) WriteFile(fn string) error {

	if fh, err := os.Create(fn); err != nil {
		return err
	} else {
		defer fh.Close()
		return this.Write(fh)
	}

	return nil
}
