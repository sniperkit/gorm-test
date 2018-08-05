package csvjson

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// WalkWokring does the work recursively on a bunch of files.
func WalkWokring(outputDir, outputExt string, delimeter rune) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var outputFilename string

		if !info.IsDir() && (filepath.Ext(path) == ".csv" || filepath.Ext(path) == ".tsv") {
			outputFilename = strings.Replace(filepath.Base(path), filepath.Ext(path), "", 1) + outputExt
			if err = Work(path, filepath.Join(outputDir, outputFilename), delimeter); err != nil {
				return err
			}
		}

		return nil
	}
}

// Work is a main wrapper function which controls the whole workflow.
func Work(input, output string, delimeter rune) (err error) {
	data, err := build(input, delimeter)
	if err != nil {
		return err
	}

	if err = save(output, data); err != nil {
		return
	}

	return nil
}

func build(input string, delimeter rune) (map[string]interface{}, error) {
	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	r.Comma = delimeter
	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return treeFrom2DMatrix(records, idsFromRecords(records), ""), nil
}

func save(output string, data map[string]interface{}) (err error) {
	out, err := os.Create(output)
	if err != nil {
		return
	}
	return json.NewEncoder(out).Encode(data)
}

// treeFrom2DMatrix is an algorithm to extract useful data from joined ID string.
// The algorightm:
// 1) keys with index 0 are all from the same node, so add them as keys and check each for final data or further run
// 2) go throuhg all keys with index 0 and find slices or maps to go further
//   1) if found a slice, then find out the amount of elements, create a slice with len(elements)
//     1) for the length of the slice run through each index and collect keys with the keyname + index
//     2) remove first elements and pass it the func, returned values added to the slice at index
func treeFrom2DMatrix(records, ids [][]string, parent string) map[string]interface{} {
	m := make(map[string]interface{})
	// fmt.Printf("ids for parent %s: \n%+v\n\n", parent, ids)

	for _, v := range ids {
		if isData(v) {
			// define a new parent
			var nparent string
			if len(parent) > 0 {
				nparent = strings.Join([]string{parent, v[0]}, "/")
			} else {
				nparent = v[0]
			}
			// get the value
			val, err := getValue(records, nparent)
			if err != nil {
				log.Fatal(err)
			}

			// assign
			valStr, ok := val.(string)
			if ok {
				valInt, err := strconv.Atoi(valStr)
				if err != nil {
					m[v[0]] = val
				} else {
					m[v[0]] = valInt
				}
			} else {
				log.Fatal("value is not convertible to string")
			}
		} else if isMap(v) {
			// define a new parent
			var nparent string
			if len(parent) > 0 {
				nparent = fmt.Sprintf("%s/%s", parent, v[0])
			} else {
				nparent = v[0]
			}

			// collect all IDs with the key
			key := cleanName(v[0])

			// create a map if it doesn't exist yet
			if _, ok := m[key]; !ok {
				cm := make(map[string]interface{})
				m[key] = cm
			}

			// collect the ids with the same key
			nids := [][]string{}
			for _, id := range ids {
				if id[0] == key {
					nids = append(nids, id[1:])
				}
			}

			// populate each key of the map with data
			if len(nids) > 0 {
				m[key] = treeFrom2DMatrix(records, nids, nparent)
			}
		} else if isSlice(v) {
			key := cleanName(v[0])

			// create a slice if it doesn't exist yet
			if _, ok := m[key]; !ok {
				// fmt.Println("slice create:", key)
				size, err := findSliceLength(idsByColumn(ids, 0), key)
				if err != nil {
					log.Fatal(err)
				}

				sl := make([]map[string]interface{}, size)
				m[key] = sl

				// define base part of the parent ID
				var nparent string
				if len(parent) > 0 {
					nparent = fmt.Sprintf("%s/%s", parent, key)
				} else {
					nparent = key
				}

				// populate each element of the slice with a map
				for i := 0; i < size; i++ {
					// fmt.Println("slice index", i)
					nids := [][]string{}
					for _, id := range ids {
						if isSlice(id) {
							// fmt.Println("slice", id)
							k, idx := getKeyIndex(id[0])
							if i == idx && k == key {
								// fmt.Println("appending", id[1:])
								nids = append(nids, id[1:]) // delete the root key
							}
						}
					}
					// fmt.Println("found", nids)
					// fmt.Println("parent", parent)
					if len(nids) > 0 {
						sl[i] = treeFrom2DMatrix(records, nids, fmt.Sprintf("%s-%d", nparent, i)) // append an index to the base part of parent
					}
				}
			}
		}
	}

	return m
}

func getValue(records [][]string, id string) (interface{}, error) {
	for _, rec := range records {
		if rec[0] == id {
			return rec[1], nil
		}
	}
	return nil, fmt.Errorf("record with id %v not found", id)
}

func idsFromRecords(records [][]string) [][]string {
	ids := make([][]string, len(records))

	for i := range records {
		ids[i] = strings.Split(records[i][0], "/")
	}

	return ids
}

func cleanName(s string) string {
	if strings.Contains(s, "-") {
		return strings.Split(s, "-")[0]
	}

	return s
}

func getKeyIndex(s string) (key string, idx int) {
	parts := strings.Split(s, "-")
	key = parts[0]
	idx, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	return
}

// isData return true if the ID has no more children and it most probably contains final data.
func isData(id []string) bool {
	if len(id) == 1 {
		return true
	}
	return false
}

// isMap returns true if IDs of a certain depth contain non-slice IDs (without '-'),
// that means that it will be a map and not a slice.
func isMap(id []string) bool {
	if len(id) > 1 {
		if !strings.Contains(id[0], "-") {
			return true
		}
	}
	return false
}

func isSlice(id []string) bool {
	if strings.Contains(id[0], "-") {
		return true
	}
	return false
}

func findSliceLength(column []string, key string) (int, error) {
	var maxIndex int

	for _, idPart := range column {
		if strings.Contains(idPart, key) {
			index := strings.Split(idPart, "-")[1]
			idx, err := strconv.Atoi(index)
			if err != nil {
				return maxIndex + 1, err
			}

			if idx > maxIndex {
				maxIndex = idx
			}
		}
	}

	// length is bigger than index by 1
	return maxIndex + 1, nil
}

// idsByColumn extracts ids by a column.
func idsByColumn(ids [][]string, colNumber int) (column []string) {
	column = []string{}

	for _, id := range ids {
		if len(id)-1 >= colNumber {
			column = append(column, id[colNumber])
		}
	}

	return
}
