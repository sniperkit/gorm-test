package main

import (
	"fmt"
	"os"

	// external
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	// internal
	"github.com/sniperkit/gorm-test/src/cmd/gorm-examples/checkExistence"
	"github.com/sniperkit/gorm-test/src/cmd/gorm-examples/compositePrimaryKey"
	"github.com/sniperkit/gorm-test/src/cmd/gorm-examples/hasMany"
)

type Example func(db *gorm.DB)

var examples = map[string]Example{
	"Has many":                  hasMany.Example,
	"Check if the record exist": checkExistence.Example,
	"Composite Primary Key":     compositePrimaryKey.Example,
}

func main() {
	for key := range examples {
		fmt.Println("-------------------------------------")
		fmt.Println(key)
		fmt.Println("-------------------------------------")
		run(key)
	}
}

func run(name string) {
	ex := examples[name]
	file := "test.db"
	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	defer os.Remove(file)
	ex(db)
	if db.Error != nil {
		fmt.Printf("Error: %s\n", db.Error)
	}
}
