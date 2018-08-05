package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/k0kubun/pp"
)

var testTables = []interface{}{
	// &Language{},
	&Tag{},
	&Star{},
	// &Repo{},
}

/*
type Topic struct {
  ID   int
  Name string
  Tag  []Tag `gorm:"polymorphic:Owner;"`
}

// Star represents a starred repository
type Repo struct {
	gorm.Model  `json:"-" yaml:"-" toml:"-"`
	RemoteID    int       `json:"id" gorm:"index;"`
	Name        *string   `gorm:"type:varchar(255);" json:"name" toml:"name" csv:"name"`
	FullName    *string   `gorm:"type:varchar(255);" json:"full_name" toml:"full_name" csv:"full_name"`
	Description *string   `gorm:"type:longtext;" json:"description" toml:"description" csv:"description"`
	Homepage    *string   `gorm:"type:varchar(255);" json:"homepage" toml:"homepage" csv:"homepage"`
	URL         *string   `gorm:"type:varchar(255);" json:"svn_url" toml:"url" csv:"repo_url"`
	Language    *string   `gorm:"type:varchar(64);" json:"language" toml:"language"  csv:"language"`
	Topics      []string  `gorm:"-" json:"topics" toml:"topics" csv:"-"`
	Stargazers  int       `json:"stargarzers_count" toml:"stargarzers_count" csv:"stargarzers_count"`
	StarredAt   time.Time `json:"starred_at" toml:"starred_at" csv:"starred_at"`
	ServiceID   uint      `json:"service_id" toml:"service_id" csv:"service_id" gorm:"index;"`
	Tags        []Tag     `gorm:"many2many:star_tags;" json:"-" toml:"-" csv:"-"`
}
*/

// Service represents a hosting service like Github
type Service struct {
	gorm.Model
	Name        string
	LastSuccess time.Time
	Stars       []Star
}

// FindOrCreateServiceByName returns a service with the specified name, creating if necessary
func FindOrCreateServiceByName(db *gorm.DB, name string) (*Service, bool, error) {
	var service Service
	if db.Debug().Where("name = ?", name).First(&service).RecordNotFound() {
		service.Name = name
		err := db.Create(&service).Error
		return &service, true, err
	}
	return &service, false, nil
}

// type User struct {
//  gorm.Model
//  Friends []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id"`
//}

// Star represents a starred repository
type Star struct {
	gorm.Model  `json:"-" yaml:"-" toml:"-"`
	RemoteID    int       `json:"id" gorm:"index;"`
	Name        *string   `gorm:"type:varchar(255);" json:"name" toml:"name" csv:"name"`
	FullName    *string   `gorm:"type:varchar(255);" json:"full_name" toml:"full_name" csv:"full_name"`
	Description *string   `gorm:"type:longtext;" json:"description" toml:"description" csv:"description"`
	Homepage    *string   `gorm:"type:varchar(255);" json:"homepage" toml:"homepage" csv:"homepage"`
	URL         *string   `gorm:"type:varchar(255);" json:"svn_url" toml:"url" csv:"repo_url"`
	Language    *string   `gorm:"type:varchar(64);" json:"language" toml:"language"  csv:"language"`
	Topics      []string  `gorm:"-" json:"topics" toml:"topics" csv:"-"`
	Stargazers  int       `json:"stargarzers_count" toml:"stargarzers_count" csv:"stargarzers_count"`
	StarredAt   time.Time `json:"starred_at" toml:"starred_at" csv:"starred_at"`
	ServiceID   uint      `json:"service_id" toml:"service_id" csv:"service_id" gorm:"index;"`
	// Tags        []Tag     `json:"-" toml:"-" csv:"-"`
	// http://gorm.io/docs/many_to_many.html#Self-Referencing
	// http://gorm.io/docs/many_to_many.html
	// http://gorm.io/docs/associations.html#Replace-Associations
	Tags []Tag `gorm:"many2many:star_tags;association_jointable_foreignkey:tag_id" json:"-" toml:"-" csv:"-"`
	// Languages []Language `gorm:"many2many:star_languages;" json:"-" toml:"-" csv:"-"` // Many-To-Many relationship, 'star_languages' is join table
	// Tags []Tag `gorm:"many2many:star_tags;foreignkey:TagId" json:"-" toml:"-" csv:"-"`
}

// CreateOrUpdateStar creates or updates a star and returns true if the star was created (vs updated)
func CreateOrUpdateStar2(db *gorm.DB, star *Star, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing Star
	if db.Where("remote_id = ? AND service_id = ?", star.RemoteID, service.ID).First(&existing).RecordNotFound() {
		star.ServiceID = service.ID
		err := db.Create(star).Error
		return err == nil, err
	}
	star.ID = existing.ID
	star.ServiceID = service.ID
	star.CreatedAt = existing.CreatedAt
	return false, db.Save(star).Error
}

func CreateOrUpdateStar(db *gorm.DB, star *Star) (*Star, bool, error) {
	// Get existing by remote ID and service ID
	var existing Star
	if db.Debug().Where("remote_id = ?", star.RemoteID).First(&existing).RecordNotFound() {
		// star.ServiceID = service.ID
		err := db.Debug().Create(star).Error
		return star, err == nil, err
	}
	star.ID = existing.ID
	// star.ServiceID = service.ID
	star.CreatedAt = existing.CreatedAt
	return star, false, db.Save(star).Error
}

func (s *Star) setTags(db *gorm.DB, topics []string, prefix string, noEmpty bool) error {

	// for _, t := range s.Tags {
	//	topics = append(topics, t.Name)
	//}

	topics = RemoveSliceDuplicates(topics, true)
	// tags := make([]Tag, 0)

	// star.Tags = append(star.Tags, *tag)
	// return db.Save(star).Error

	pp.Println(topics)

	for _, topic := range topics {
		name := fmt.Sprintf("%s%s", prefix, topic)
		if len(strings.TrimSpace(name)) == 0 {
			continue
		}
		tag := Tag{Name: name}
		// var existing Tag
		if db.Debug().Where("lower(name) = ?", strings.ToLower(name)).First(&tag).RecordNotFound() {
			tag.Name = name
			if err := db.Debug().Create(&tag).Error; err != nil {
				log.Fatalln("error while creating new tag...")
			}
		}
		/*
			err := db.Debug().FirstOrCreate(&tagModel, tag).Error
			if err != nil {
				log.Fatalln("FirstOrCreate.Error: ", err)
				return err
			}
		*/
		// tags = append(tags, tag)
		s.Tags = append(s.Tags, tag)
	}
	// s.Tags = append(s.Tags, tags...)
	// s.Tags = tags
	// return false, db.Save(s).Error
	return db.Debug().Save(&s).Error
}

// FindOrCreateTagByName finds a tag by name, creating if it doesn't exist
func FindOrCreateTagByName(db *gorm.DB, name string) (*Tag, bool, error) {
	var tag Tag
	if db.Debug().Where("lower(name) = ?", strings.ToLower(name)).First(&tag).RecordNotFound() {
		tag.Name = name
		err := db.Debug().Create(&tag).Error
		return &tag, true, err
	}
	return &tag, false, nil
}

// AddTag adds a tag to a star
func (star *Star) AddTag(db *gorm.DB, tag *Tag) error {
	star.Tags = append(star.Tags, *tag)
	return db.Save(star).Error
}

func SetTags(db *gorm.DB, topics []string, prefix string, noEmpty bool) error {
	topics = RemoveSliceDuplicates(topics, true)
	tags := make([]Tag, 0)
	for _, topic := range topics {
		name := fmt.Sprintf("%s%s", prefix, topic)
		if len(strings.TrimSpace(name)) == 0 {
			continue
		}
		tag := Tag{Name: fmt.Sprintf("%s%s", prefix, topic)}
		var tagModel Tag
		err := db.Debug().FirstOrCreate(&tagModel, tag).Error
		if err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	return nil
}

type Tag2 struct {
	// ID int
	gorm.Model
	// UserID  int     `gorm:"index"` // Foreign key (belongs to), tag `index` will create index for this column
	Name string `gorm:"type:varchar(100);unique_index"` // `type` set sql type, `unique_index` will create unique index for this column
	// Enabled bool
	Stars []Star `gorm:"many2many:star_tags;" json:"-" yaml:"-" toml:"-"`
}

type Tag struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	// gorm.Model `json:"-" yaml:"-" toml:"-"`
	Name      string `gorm:"unique;not null" json:"name" yaml:"name" toml:"name" csv:"name"`
	StarCount int    `gorm:"-" json:"star_count" yaml:"star_count" toml:"star_count" toml:"csv"`
	Stars     []Star `gorm:"many2many:star_tags;" json:"-" yaml:"-" toml:"-"`
}

type Language struct {
	ID   int
	Name string `gorm:"index:idx_name_code"` // Create index with name, and will create combined index if find other fields defined same name
	Code string `gorm:"index:idx_name_code"` // `unique_index` also works
}

func truncateTables(db *gorm.DB, tables ...interface{}) {

	for _, table := range tables {
		if err := db.Debug().DropTableIfExists(table).Error; err != nil {
			log.Fatalln("error while dropping table: ", err)
		}
		if err := db.Debug().DropTableIfExists(table).Error; err != nil {
			log.Fatalln("error while dropping table: ", err)
		}
		db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")
		if err := db.Debug().AutoMigrate(table).Error; err != nil {
			log.Fatalln("error while auto-migrating table: ", err)
		}
	}
}
