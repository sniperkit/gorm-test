package main

import (
	"fmt"

	// external
	"github.com/jinzhu/gorm"
)

type UserMany struct {
	ID    uint `gorm:"primary_key"`
	Name  string
	Roles []Role
}

type Role struct {
	Name   string
	UserID uint
}

func hasManyExample(db *gorm.DB) {
	// init db
	db.AutoMigrate(&UserMany{})
	db.AutoMigrate(&Role{})
	db.Model(&UserMany{}).Related(&Role{})
	// create user1 in 'users' table and role1, role2 in 'roles' table.
	user := UserMany{
		Name: "user1",
		Roles: []Role{
			{Name: "role1"},
			{Name: "role2"},
		},
	}
	db.Create(&user)

	// find the created user.  for eager loading, 'Preload("Roles")' is required.
	found := &UserMany{}
	db.Preload("Roles").Find(&found)
	fmt.Println(found)
}
