package main

import (
	"fmt"

	// external
	"github.com/jinzhu/gorm"
)

type UserCPK struct {
	Name  string `gorm:"primary_key"`
	Email string `gorm:"primary_key"`
}

type MutedLogger struct {
}

func (m MutedLogger) Print(v ...interface{}) {
	// do nothing
}

func compositePrimaryKeyExample(db *gorm.DB) {
	db.AutoMigrate(&UserCPK{})
	user1 := UserCPK{
		Name:  "a",
		Email: "a@foo.com",
	}
	if err := db.Create(&user1).Error; err != nil {
		panic(err)
	}

	user2 := UserCPK{
		Name:  "a",
		Email: "a@foo.com",
	}
	db.SetLogger(MutedLogger{})
	if err := db.Create(&user2).Error; err != nil {
		// expected behavior
	} else {
		panic("duplicated record must not be inserted")
	}
	db.SetLogger(gorm.Logger{})

	user3 := UserCPK{
		Name:  "a",
		Email: "a@bar.com",
	}
	if err := db.Create(&user3).Error; err != nil {
		panic(err)
	} else {
		// expected behavior
	}

	allUsers := []UserCPK{}
	db.Find(&allUsers)
	fmt.Println("Created users:")
	for _, v := range allUsers {
		fmt.Printf("%#v\n", v)
	}
}
