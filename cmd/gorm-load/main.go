package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/k0kubun/pp"
)

var (
	localDumpFile string = "./gh-starred.json"
	withDialect   string = "mysql"
	withIndexer   string = "manticore" // elasticsearch or manticore/sphinxsearch
)

/*
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

func main() {

	// db.Set("gorm:insert_option", "ON DUPLICATE KEY UPDATE").Create(&star)
	// ON DUPLICATE KEY UPDATE

	var db *gorm.DB
	var err error
	switch withDialect {
	case "mysql":
		db, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", "gorm_test", "gorm1234T!", "127.0.0.1", "3306", "gorm_test"))
		if err != nil {
			log.Fatalln("error while creating connection with database: ", err)
		}
	case "sqlite", "sqlite3":
		fallthrough
	default:
		db, err = gorm.Open("sqlite", "gh_starred.db")
		if err != nil {
			log.Fatalln("error while creating connection with database: ", err)
		}
	}

	defer db.Close()
	db.LogMode(true)

	truncateTables(db, testTables...)

	content, err := ioutil.ReadFile(localDumpFile)
	if err != nil {
		log.Fatalln("error while reading file: ", err)
	}

	var starred []Star
	if err = json.Unmarshal(content, &starred); err != nil {
		log.Fatalln("error while Unmarshaling: ", err)
	}

	var i = 0
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

		// if err := db.Model(&star).Association("Tags").Error; err != nil {
		//	log.Println("error while creating association: ", err)
		//}

		pp.Println("FULLNAME=", *star.FullName, ", URI=", uri)

		// for _, t := range dtags {
		// star.setTags(db, dtags, "", true)
		//}

		if err := star.setTags(db, dtags, "snk/", true); err != nil {
			log.Fatalln("error while saving star into db: ", err)
		}

		if err := db.Save(&star).Error; err != nil {
			log.Println("error while saving star into db: ", err)
		}

		// if err := db.Set("gorm:insert_option", "ON DUPLICATE KEY UPDATE").Model(&star).Association("Tags").Append(Topics2Tags(dtags, "", true)).Error; err != nil {
		//if err := db.Model(&star).Association("Tags").Append(Topics2Tags(dtags, "", true)).Error; err != nil {
		//	log.Println("error while appending to stars/tags association: ", err)
		//}

		// err := db.Save(star).Association("Tags").Error
		// if err := db.Save(&star).Error; err != nil {
		//	log.Println("error while saving data into db: ", err)
		// }
		// pp.Println(star)

		i++

	}

	pp.Println("processed: ", i)

}
