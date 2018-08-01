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

func main() {

	// db.Set("gorm:insert_option", "ON DUPLICATE KEY UPDATE").Create(&star)
	// ON DUPLICATE KEY UPDATE

	var db *gorm.DB
	var err error
	switch withDialect {
	case "mysql":
		db, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", "gorm_test", "gorm1234!T", "127.0.0.1", "3306", "gorm_test"))
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

		// topics := strings.Join(star.Topics, ", ")
		aggr := fmt.Sprintf(`%s %s %s %s %s`, uri, desc, readme, lang, strings.Join(star.Topics, ", "))
		stacks := ExtractStack(aggr, uri)
		stacks = AddTopicsPrefix(stacks, "stackexchange/", true)
		ghTopics := AddTopicsPrefix(star.Topics, "github/", true)

		var dtags []string
		dtags = append(dtags, stacks...)
		dtags = append(dtags, ghTopics...)
		dtags = RemoveSliceDuplicates(dtags, true)
		dtags = AddTopicsPrefix(dtags, "", true)

		// star.Tags = Topics2Tags(dtags, "", true)
		// star.DT = Topics2Tags(dtags, "", true)
		// star.DT = dtags

		// if err := db.Model(&star).Association("Tags").Error; err != nil {
		//	log.Println("error while creating association: ", err)
		//}

		pp.Println("FULLNAME=", *star.FullName, ", URI=", uri)

		// if err := db.Set("gorm:insert_option", "ON DUPLICATE KEY UPDATE").Model(&star).Association("Tags").Append(Topics2Tags(dtags, "", true)).Error; err != nil {
		if err := db.Model(&star).Association("Tags").Append(Topics2Tags(dtags, "", true)).Error; err != nil {
			log.Println("error while appending to stars/tags association: ", err)
		}

		// err := db.Save(star).Association("Tags").Error
		if err := db.Save(&star).Error; err != nil {
			log.Println("error while saving data into db: ", err)
		}
		// pp.Println(star)
		i++
	}

	pp.Println("processed: ", i)

}
