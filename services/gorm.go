package services

import (
	"consulate/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

var _gormDb *gorm.DB
var _gormDbCreator sync.Once

func GormDb() *gorm.DB {
	_gormDbCreator.Do(func() {
		db, err := gorm.Open(sqlite.Open(GetConfig().Database.Path), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		err = db.AutoMigrate(&models.Enquiry{}, &models.FollowUp{})
		if err != nil {
			panic(err)
		}

		_gormDb = db
	})

	return _gormDb
}
