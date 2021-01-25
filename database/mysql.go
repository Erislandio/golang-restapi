package database

import (
	"database/sql"
	"fmt"

	"github.com/erislandio/web/restapi/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

// Init database mysql
func Init() {

	conf := config.GetConfig()

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.DB_USER, conf.DB_PASS, conf.DB_HOST, conf.DB_PORT, conf.DB_NAME)

	db, err = sql.Open("mysql", connectionString)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database connected with success!")

}

//  CreateMysqlConn func retunr mysql database con
func CreateMysqlConn() *sql.DB {
	return db
}
