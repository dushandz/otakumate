package comicdao

import (
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ComicTable struct {
	ID        int64 `gorm:"primary_key"`
	Name      string
	Cover     string
	TP        string
	CreatedAt time.Time
}
type VolsTable struct {
	ID        int64 `gorm:"primary_key"`
	ComicID   int64
	Name      string
	Vols      string
	CreatedAt time.Time
}

type VolTable struct {
	ID        int64 `gorm:"primary_key;AUTO_INCREMENT"`
	VolID     int64
	Image     string
	CreatedAt time.Time
}

type ComicDao struct {
	db *gorm.DB
}

func NewComicDao() *ComicDao {
	conectDB()
	var dao = &ComicDao{db: db}
	return dao
}

func (d ComicDao) InsertComicData(comic *ComicTable) {
	if !d.db.HasTable(&ComicTable{}) {
		d.db.CreateTable(&ComicTable{})
	}
	if d.db.NewRecord(comic) {
		d.db.Create(comic)
	} else {
		d.db.Save(&comic)
	}
}

func (d ComicDao) QueryComic(tp string) []ComicTable {
	list := make([]ComicTable, 0)
	if !d.db.HasTable(&ComicTable{}) {
		d.db.CreateTable(&ComicTable{})
	} else {
		d.db.Where("TP = ?", tp).Find(&list)
	}
	return list
}

func (d ComicDao) InsertVolsData(vols *VolsTable) {
	if !d.db.HasTable(&VolsTable{}) {
		d.db.CreateTable(&VolsTable{})
	}
	if d.db.NewRecord(vols) {
		d.db.Create(vols)
	} else {
		d.db.Save(&vols)
	}
}

func (d ComicDao) QueryVols(volid int64) []VolsTable {
	list := make([]VolsTable, 0)
	if !d.db.HasTable(&VolsTable{}) {
		d.db.CreateTable(&VolsTable{})
	} else {
		d.db.Where("ID = ?", volid).Find(&list)
	}
	return list
}

func (d ComicDao) InsertVolDetailData(vol *VolTable) {
	if !d.db.HasTable(&VolTable{}) {
		d.db.CreateTable(&VolTable{})
	}
	if d.db.NewRecord(vol) {
		d.db.Create(vol)
	} else {
		d.db.Save(&vol)
	}
}

func (d ComicDao) QueryVolDetail(volid int64) []VolTable {
	list := make([]VolTable, 0)
	if !d.db.HasTable(&VolTable{}) {
		d.db.CreateTable(&VolTable{})
	} else {
		d.db.Where("VolID = ?", volid).Find(&list)
	}

	return list
}
