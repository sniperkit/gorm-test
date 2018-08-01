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

func main() {

	// db.Set("gorm:insert_option", "ON DUPLICATE KEY UPDATE").Create(&star)
	// ON DUPLICATE KEY UPDATE

	db, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", "root", "da33ve79T!", "127.0.0.1", "3306", "snk_gorm_test"))
	if err != nil {
		log.Fatalln("error while creating connection with database: ", err)
	}

	defer db.Close()
	db.LogMode(true)

	// db.AutoMigrate(&Star{})
	// db.AutoMigrate(&Tag{})
	truncateTables(db, testTables...)

	// db.AutoMigrate(&Topic{})
	content, err := ioutil.ReadFile("gh-starred.json")
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
