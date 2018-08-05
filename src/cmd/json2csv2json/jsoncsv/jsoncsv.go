package jsoncsv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// WalkWokring does the work recursively on a bunch of files.
func WalkWokring(outputDir, outputExt string, delimeter rune) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var outputFilename string

		if !info.IsDir() && filepath.Ext(path) == ".json" {
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
	data := make(map[string]interface{})
	if err = build(input, data); err != nil {
		return
	}

	if err = save(output, delimeter, data); err != nil {
		return
	}

	return nil
}

// build reads a source file and reflects on its content to build a plain map of id:value pairs.
func build(filename string, output map[string]interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data := make(map[string]interface{})

	d := json.NewDecoder(f)
	if err = d.Decode(&data); err != nil {
		return err
	}

	// convert to plain id:value map
	for k, v := range data {
		iterate(k, v, output)
	}

	return nil
}

// iterate goes recursively through an interface value, composes ids and values and adds them to the output map.
func iterate(parent string, data interface{}, output map[string]interface{}) {
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.String:
		// fmt.Printf("%s (string)\n", parent)
		output[parent] = v.Interface()
		return
	case reflect.Float64:
		// fmt.Printf("%s (float64)\n", parent)
		output[parent] = v.Interface()
		return
	case reflect.Slice:
		// fmt.Printf("%s (slice of %v)\n", parent, v.Len())
		for i := 0; i < v.Len(); i++ {
			iterate(
				fmt.Sprintf("%s-%d", parent, i),
				reflect.ValueOf(v.Index(i).Interface()).Interface(),
				output,
			)
		}
		return
	case reflect.Map:
		// fmt.Printf("%s (map with %v)\n", parent, v.MapKeys())
		for _, k := range v.MapKeys() {
			iterate(
				parent+"/"+k.String(),
				v.MapIndex(k).Interface(),
				output,
			)
		}
		return
	}
}

// save saves an output map to a CSV-, TSV-file. The format of the file is dependent on the delimeter
// provided to the function: "\t" or ",".
// TODO: use encoding/csv
func save(filename string, delimeter rune, output map[string]interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = delimeter

	defer w.Flush()

	if err = w.Write([]string{"#id", "value"}); err != nil {
		return err
	} // header

	for id, value := range output {
		// NOTE: specific requirements, can be ignored for more general usage
		// write only fields: title, text, lead, teaser, description
		// if strings.Contains(id, "title") || strings.Contains(id, "text") || strings.Contains(id, "lead") || strings.Contains(id, "teaser") || strings.Contains(id, "description") {
		// 	if err = w.Write([]string{id, fmt.Sprintf("%v", value)}); err != nil {
		// 		return err
		// 	} // rows
		// }

		var v interface{}
		vStr, ok := value.(string)
		if ok {
			v = html.UnescapeString(vStr)
		} else {
			vInt, ok := value.(float64)
			if ok {
				v = vInt
			} else {
				return fmt.Errorf("value is unconvertible to string and int: %v", value)
			}
		}

		if err = w.Write([]string{id, fmt.Sprintf("%v", v)}); err != nil {
			return err
		} // rows
	}

	return nil
}

//

// func idsTreeFromRecords(root map[string]interface{}, records [][]string) {
// 	var maxIndex int

// 	for _, rec := range records {
// 		// split the whole id for parts
// 		parts := strings.Split(rec[0], "/")

// 		// fill out end-nodes
// 		if len(parts) == 1 {
// 			root[parts[0]] = "end-node"
// 			continue
// 		}

// 		// get key and index for the first part
// 		k, idx := getKeyIndex(parts[0])
// 		if idx > maxIndex {
// 			// update the max index
// 			maxIndex = idx
// 		}

// 		if _, ok := root[k]; !ok {
// 			m := make(map[int]interface{})
// 			root[k] = m

// 			if strings.Contains(parts[1], "-") {
// 				m[idx] = idsTreeFromRecords(m, records)
// 			}
// 		} else {
// 			if m, ok := root[k].(map[int]interface{}); ok {
// 				m[idx] = parts[1:]
// 			}
// 		}
// 	}

// 	fmt.Printf("%+v\n", root)
// }

// func mapFromIdsByColumn(ids [][]string) (map[string]interface{}, error) {
// 	m := make(map[string]interface{})

// 	for _, id := range ids {
// 		println(id)

// 	}

// 	return m, nil
// }

// func mapFromIds(ids [][]string) map[string]interface{} {
// 	m := make(map[string]interface{})

// 	for _, v := range ids {
// 		// updateMapFromID(m, v, ids)
// 	}

// 	return m
// }

// func mapFromID(id []string, ids [][]string, depth int) map[string]interface{} {
// 	m := make(map[string]interface{})

// 	if depth > len(id)-1 {
// 		return m
// 	}

// 	if isData(id, depth) {
// 		m[id[depth]] = "data"
// 		return m
// 	}

// 	if isMap(ids, depth) {
// 		keys := keysFromColumn(idsByColumn(ids, depth))
// 		for k, t := range keys {
// 			switch t {
// 			case "slice":
// 				size, err := findSliceLength(idsByColumn(ids, depth+1), k)
// 				if err != nil {
// 					log.Fatal(err)
// 				}

// 				fmt.Println("size", size, idsByColumn(ids, depth+1))
// 				slice := make([]map[string]interface{}, size)
// 				m[k] = slice

// 				_, idx := getKeyIndex(id[depth+1])
// 				fmt.Println("idx", idx)
// 				slice[idx] = mapFromID(id, ids, depth+1)
// 			case "data or map":
// 				m[k] = t
// 			}
// 		}
// 	}

// 	return m
// }

// func iterByColumn(ids [][]string, depth int) map[string]interface{} {
// 	m := make(map[string]interface{})

// 	if isMap(ids, depth) {
// 		keys := keysFromColumn(idsByColumn(ids, depth))
// 		fmt.Printf("keys: %+v depth: %v\n", keys, depth)

// 		for k, t := range keys {
// 			switch {
// 			case t == "data":
// 				m[k] = t
// 				continue
// 			case strings.Contains(t, "slice"):
// 				curCol := idsByColumn(ids, depth)

// 				size, err := findSliceLength(curCol, k)
// 				if err != nil {
// 					log.Fatal(err)
// 				}

// 				slice := make([]map[string]interface{}, size)
// 				m[k] = slice

// 				// fmt.Printf("curCol: %v\n", curCol)
// 				// go through keys with full names
// 				// for _, fullKey := range curCol {
// 				// 	if strings.Contains(fullKey, k) {
// 				// 		key, curIdx := getKeyIndex(fullKey)
// 				// 		fmt.Printf("curIdx %s %d\n", key, curIdx)
// 				// 		slice[curIdx] = iterByColumn(ids, depth+1)
// 				// 		// return m
// 				// 	}
// 				// }

// 				for _, id := range ids {
// 					if len(id)-1 < depth {
// 						continue
// 					}

// 					fmt.Println("id", id)

// 					for i := 0; i < size; i++ {
// 						if id[depth] == fmt.Sprintf("%s-%d", k, i) {
// 							slice[i] =
// 						}
// 					}
// 				}

// 				// fmt.Printf("slice: %+v\n", slice)
// 			}
// 		}
// 	}

// 	return m
// }

// func defineType(ids [][]string, id []string, m map[string]interface{}) {
// 	// run in depth
// 	for depth := 0; depth < len(id); depth++ {
// 		defineTypeIter(ids, id, depth)
// 	}
// }

// func defineType(ids [][]string, id []string, m map[string]interface{}) {
// 	// run in depth
// 	defineTypeIter(ids, id, 0)
// }

// func defineTypeIter(ids [][]string, id []string, m map[string]interface{}, depth int) {
// 	// data
// 	if len(id)-1 == depth {
// 		m[id[depth]] = "data"
// 		return
// 	}

// 	idParts := strings.Split(id[depth], "-")

// 	if len(idParts) > 1 {
// 		// slice
// 		size, err := findSliceLength(idsByColumn(ids, 0), idParts[0])
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// slice of maps or other slices?
// 		if _, ok := m[idParts[0]]; !ok {
// 			m[idParts[0]] = make([]map[string]interface{}, size)
// 		}

// 		for i := 0; i < size; i++ {
// 			defineTypeIter(ids, id, m[idParts[0]][i], depth+1)
// 		}
// 	} else {
// 		// map
// 		if _, ok := m[id[depth]]; ok {
// 			log.Println("map already exists (error)")
// 		} else {
// 			m[id[depth]] = make(map[string]interface{})
// 		}
// 	}
// }

// func idsByColumnRecursively(ids [][]string, level int, m map[int][]string) {
// 	var more bool
// 	col := []string{}

// 	for _, id := range ids {
// 		if len(id) >= level {
// 			col = append(col, id[0])
// 		}

// 		if len(id) > level {
// 			more = true
// 		}
// 	}

// 	m[level] = col

// 	if more {
// 		idsByColumnRecursively(ids, level+1, m)
// 	}
// }

// func Analyze(id string, m map[string]interface{}) {
// 	parts := strings.Split(id, "/")

// 	if len(parts) == 1 {
// 		m[parts[0]] = "data"
// 		return
// 	}

// 	parts1 := strings.Split(parts[0], "-")
// 	if len(parts1) > 1 {
// 		m[parts1[0]] = []map[string]interface{}{}
// 		Analyze(strings.Join(parts[1:], "/"), m[parts1[0]])
// 	} else {
// 		m[parts[0]] = parts[1:]
// 	}
// }
