package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	// "strconv"
	"strings"

	// external
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	// "github.com/fatih/structs"
	// "github.com/gocarina/gocsv"
	"github.com/jinzhu/configor"
	"github.com/k0kubun/pp"
)

var (
	localDumpFile string = "./shared/testdata/gh-starred.json"
	// localDumpFile string = "../../shared/testdata/gh-starred.json"
	withDialect string = "mysql"
	withIndexer string = "manticore" // elasticsearch or manticore/sphinxsearch
)

/*
	- /usr/local/opt/mysql@5.7/bin/mysql.server start

	Refs:
	- disable the NO_ZERO_DATE in your sql mode.
		- select @@GLOBAL.sql_mode;
		- SET GLOBAL sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
	- LOAD DATA INFILE
		- https://github.com/swgloomy/askQuestions/blob/master/main.go
		- https://github.com/percona/qan-api/blob/master/tests/setup/db/db.go#L134
		- https://github.com/mattermost/mattermost-load-test/blob/master/loadtest/bulkload.go#L715
	- ASSOCIATIONS
		- https://github.com/ranjur/go-gin-example/blob/master/articles/models.go#L39-L43
	- INSIGHTS/FEEDBACKS
		- https://github.com/Depado/articles/blob/master/pages/gorm-gotchas.md
*/

func loadConfig() {
	configor.Load(cfg, "config.yml")
	pp.Printf("config: %#v", cfg)
}

func main() {

	// db.Set("gorm:insert_option", "ON DUPLICATE KEY UPDATE").Create(&star)
	// ON DUPLICATE KEY UPDATE

	log.Println("withDialect: ", withDialect)

	var db *gorm.DB
	var err error
	switch withDialect {
	case "postgres":
		fallthrough
	case "mssql":
		log.Fatalln("database type ", withDialect, "not ready yet...")
	case "mariadb":
		fallthrough
	case "mysql":
		db, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local&charset=utf8mb4", "root", "da33ve79T!", "127.0.0.1", "3306", "snk_gorm_test"))
		// db, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", "root", "gorm_super_secret", "mariadb", "3306", "gorm_mariadb_test"))
		if err != nil {
			log.Fatalln("error while creating connection with database (", withDialect, "): ", err)
		}
	case "sqlite", "sqlite3":
		fallthrough
	default:
		db, err = gorm.Open("sqlite", "/data/gh_starred.db")
		if err != nil {
			log.Fatalln("error while creating connection with database: ", err)
		}
	}

	defer db.Close()
	db.LogMode(true)
	// db.Set("gorm:table_options", "charset=utf8")
	// db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")

	// db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 auto_increment=1")
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 auto_increment=1")
	truncateTables(db, testTables...)

	content, err := ioutil.ReadFile(localDumpFile)
	if err != nil {
		log.Fatalln("error while reading file: ", err)
	}

	var starred []Star
	if err = json.Unmarshal(content, &starred); err != nil {
		log.Fatalln("error while Unmarshaling: ", err)
	}

	/*
		starredWriterCSV, err := newCsvWriter("star.csv")
		if err != nil {
			log.Fatalln("error while creating csv writer for starred repos: ", err)
		}

		topicsWriterCSV, err := newCsvWriter("tags.csv")
		if err != nil {
			log.Fatalln("error while creating csv writer for starred topics: ", err)
		}
	*/

	var i = 0

	// var stars []Star

	// db.Model(&user).Association("Languages")

	for _, star := range starred {
		var uri, desc, readme, lang string

		uri = fmt.Sprintf("github.com/%s", *star.FullName)

		if star.Description != nil {
			desc = *star.Description
		}

		if star.Language != nil {
			lang = *star.Language
		}

		aggr := fmt.Sprintf(`%s %s %s %s %s`, uri, desc, readme, lang, strings.Join(star.Topics, ", "))
		stacks := ExtractStack(aggr, uri)
		stacks = AddTopicsPrefix(stacks, "stackexchange/", true)
		ghTopics := AddTopicsPrefix(star.Topics, "github/", true)

		var dtags []string
		dtags = append(dtags, stacks...)
		dtags = append(dtags, ghTopics...)
		dtags = RemoveSliceDuplicates(dtags, true)
		dtags = AddTopicsPrefix(dtags, "", true)

		pp.Println("FULLNAME=", *star.FullName, ", URI=", uri)

		// star.Tags = Topics2Tags(dtags, "", true)
		if _, _, err := CreateOrUpdateStar(db, &star); err != nil {
			log.Fatalln("CreateOrUpdateStar.Error: ", err)
		}

		star.setTags(db, dtags, "", true)

		i++

	}

	pp.Println("processed: ", i)

}

func GetFields(i interface{}) (res []string) {
	v := reflect.ValueOf(i)
	for j := 0; j < v.NumField(); j++ {
		res = append(res, v.Field(j).String())
	}
	return
}
