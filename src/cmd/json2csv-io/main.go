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

package main

import (
	"flag"
	"github.com/jscherff/json2csv/io"
	"log"
	"os"
)

const dstFile string = `out.csv`

var (
	fSrcUrl  = flag.String(`source-url`, ``, "Read JSON objects from `<url>`")
	fSrcFile = flag.String(`source-file`, ``, "Read JSON objects from `<file>`")
	fDstFile = flag.String(`dest-file`, dstFile, "Write CSV records to `<file>`")
)

func init() {

	log.SetFlags(0)
	flag.Parse()

	if *fSrcUrl == `` && *fSrcFile == `` {
		log.Print(`You must specify a source file or URL`)
		flag.Usage()
		os.Exit(1)
	}
}

func main() {

	var err error

	rw := io.NewReadWriter()

	switch true {
	case *fSrcUrl != ``:
		err = rw.ReadUrl(*fSrcUrl)
	case *fSrcFile != ``:
		err = rw.ReadFile(*fSrcFile)
	}

	if err == nil {
		err = rw.WriteFile(*fDstFile)
	}

	if err != nil {
		log.Fatal(err)
	}
}
