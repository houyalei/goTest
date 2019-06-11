package models

import (
	"log"
	"work/pkg/setting"

	"github.com/kirinlabs/mysqldb"
)

var db *mysqldb.Adapter
var err error

func init() {
	var (
		err                          error
		dbName, user, password, host string
	)
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database':%v", err)
	}
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()

	db, err = mysqldb.New(
		&mysqldb.Options{
			User:         user,
			Password:     password,
			Host:         host,
			Port:         3306,
			Database:     dbName,
			Charset:      "utf8",
			MaxIdleConns: 5,
			MaxOpenConns: 10,
			Debug:        true,
		})
	if err != nil {
		log.Panic("Connect mysql server error: ", err)
	}
}
