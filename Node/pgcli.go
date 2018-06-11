package spider

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func init() {
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=gorm sslmode=disable password=68957423")
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&Node{}) //自迁移模式
}
