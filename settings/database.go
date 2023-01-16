package settings

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"partitions_searcher/proto/partitions_searcher"
)

var Db *gorm.DB

func Init(configPath string) {
	configs := GetConf(configPath)
	initDB(configs.DbHost, configs.DbUser, configs.DbPassword, configs.DbName, configs.DbPort)
}

func initDB(host string, user string, password string, dbName string, port string) {
	var err error
	dbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	Db, err = gorm.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	Db.AutoMigrate()
}

func GetEventByPartition(partition string, query string) (records []partitions_searcher.Record, err error) {
	err = Db.Raw(fmt.Sprintf("select id, title, description from %s as e%s", partition, query)).Scan(&records).Error
	return
}
