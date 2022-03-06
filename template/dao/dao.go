package dao

import (
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB
}

func GetDaoInstance() *Dao {
	return newDao(dbclt.Engine)
}
func newDao(DbEngine *gorm.DB) *Dao {
	return &Dao{db: DbEngine.Session(&gorm.Session{})}
}
