package model

import (
	"github.com/Unknwon/goconfig"
	log "github.com/go-irain/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var DBMeeting *gorm.DB
var DBCat *gorm.DB

// InitDB ...
func InitDB(config *goconfig.ConfigFile) {
	var err error
	DB, err = gorm.Open("mysql", config.MustValue("db", "race"))
	if err != nil {
		log.Debug("error in open db")
		// 首次连接数据库失败，直接退出进程
		panic(err.Error())
	}

	// 设置连接池
	DB.DB().SetMaxIdleConns(5)
	DB.DB().SetMaxOpenConns(10)
	// 全局禁用表名复数
	DB.SingularTable(true)

	// DB.LogMode(true)
	// DB.SetLogger(gorm.Logger{revel.TRACE})

	log.Debug("success connect to DB:race")

	// 进程运行期间保持，不关闭数据库连接
	// DB.Close()

	DBMeeting, err = gorm.Open("mysql", config.MustValue("db", "meeting"))
	if err != nil {
		log.Debug(0, "error in open db")
		// 首次连接数据库失败，直接退出进程
		panic(err.Error())
	}

	// 设置连接池
	DBMeeting.DB().SetMaxIdleConns(5)
	DBMeeting.DB().SetMaxOpenConns(10)
	// 全局禁用表名复数
	DBMeeting.SingularTable(true)
	// DBI2P.LogMode(true)

	DBCat, err = gorm.Open("mysql", config.MustValue("db", "category"))
	if err != nil {
		log.Debug(0, "error in open db")
		// 首次连接数据库失败，直接退出进程
		panic(err.Error())
	}

	// 设置连接池
	DBCat.DB().SetMaxIdleConns(5)
	DBCat.DB().SetMaxOpenConns(10)
	// 全局禁用表名复数
	DBCat.SingularTable(true)
	// DBCONFIG.LogMode(true)
	log.Debug("success connect to DB:DBCat")
}

// TabRace ...
type TabRace struct {
	RaceID          string `gorm:"column:race_id" json:"race_id"`
	RaceName        string `gorm:"column:race_name" json:"race_name"`
	RaceNumber      int    `gorm:"column:race_number" json:"race_number"`
	AdvertisedStart string `gorm:"column:advertised_start" json:"advertised_start"`
	MeetingID       string `gorm:"column:meeting_id" json:"meeting_id"`
}

// TableName ...
func (TabRace) TableName() string {
	return "tab_race"
}

// GetRace ...
func GetRace(page int, pagesize int) (data []TabRace, err error) {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	err = DB.Model(&TabRace{}).
		// Where("to_id=?", kefuId).
		Offset(offset).
		Limit(pagesize).
		// Order("status desc, updated_at desc").
		Find(&data).Error
	return data, err
}
