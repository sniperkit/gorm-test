package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

var testTables = []interface{}{
	&Star{},
	&Topic{},
}

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`
}

type User struct {
	Model
	Email string `validator:"email" gorm:"unique_index" json:"email"`
}

type Category struct {
	Model
	Name string `json:"name"`
}

type Topic struct {
	Model
	CategoryID uint
	Category   Category
	Title      string `json:"title"`
	Acl        []User `json:"acl" gorm:"many2many:topic_acl"`
}

// Star represents a starred repository
type Star struct {
	gorm.Model
	RemoteID    string   // `gorm:"index;"`
	Name        *string  `gorm:"type:varchar(255);" json:"name"`
	FullName    *string  `gorm:"type:varchar(255);" json:"full_name"`
	Description *string  `gorm:"type:longtext;" json:"description"`
	Homepage    *string  `gorm:"type:varchar(255);" json:"homepage"`
	URL         *string  `gorm:"type:varchar(255);" json:"svn_url"`
	Language    *string  `gorm:"type:varchar(64);" json:"language"`
	Topics      []string `gorm:"-" json:"topics"`
	Stargazers  int      `json:"stargarzers_count"`
	StarredAt   time.Time
	ServiceID   uint  // `gorm:"index;"`
	Tags        []Tag `gorm:"many2many:star_tags;"`
	DT          []Tag `gorm:"-"`
}

type Tag struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255);unique_index;not null;"` // `type:varchar(255)`
	Stars     []Star `gorm:"many2many:star_tags;"`
	StarCount int    `gorm:"-"`
}

func truncateTables(db *gorm.DB, tables ...interface{}) {
	for _, table := range tables {
		if err := db.DropTableIfExists(table).Error; err != nil {
			log.Fatalln("error while dropping table: ", err)
		}
		if err := db.AutoMigrate(table).Error; err != nil {
			log.Fatalln("error while auto-migrating table: ", err)
		}
	}
}
