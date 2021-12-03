package spider

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

func init() {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&Node{}) //自迁移模式
}
