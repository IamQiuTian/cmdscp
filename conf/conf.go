package conf

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

type People struct {
	ip       string
	username string
	password string
	port     int
}

func ConnectDB(driverName string, dbName string) *sql.DB {
	db, err := sql.Open(driverName, dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func Read(c *sql.DB, ip string) (InfoMap map[string]interface{}) {
	InfoMap = make(map[string]interface{})
	rows, err := c.Query("select * from demo")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		p := new(People)
		err := rows.Scan(&p.ip, &p.port, &p.username, &p.password)
		if err != nil {
			log.Fatal(err)
		}
		if ip != strings.Replace(p.ip, " ", "", -1) {
			continue
		}

		InfoMap["ip"] = strings.Replace(p.ip, " ", "", -1)
		InfoMap["port"] = p.port
		InfoMap["username"] = strings.Replace(p.username, " ", "", -1)
		InfoMap["password"] = strings.Replace(p.password, " ", "", -1)
	}
	return
}

func FindInfo(ip string) (InfoMap map[string]interface{}) {
	c := ConnectDB("sqlite3", "conf/info.db")
	InfoMap = Read(c, ip)
	return
}
