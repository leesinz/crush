package database

import (
	"crush/config"
	"crush/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var packetstormdb *sql.DB

func init() {
	conf := config.LoadConfig()
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s", conf.Database.DBUsername, conf.Database.DBPassword, conf.Database.DBPort, conf.Database.Name)
	var err error
	packetstormdb, err = sql.Open("mysql", connectionString)
	if err != nil {
		utils.PrintColor("error", "Error connecting mysql database:", err, "\n")
		return
	}
	packetstormdb.SetMaxIdleConns(20)
	packetstormdb.SetMaxOpenConns(50)

	if err := packetstormdb.Ping(); err != nil {
		utils.PrintColor("error", "Error connecting mysql database:", err, "\n")
		return
	}
}

func CreatePacketstormDB() {
	_, err := packetstormdb.Exec(`
		CREATE TABLE IF NOT EXISTS packetstorm_db (
			ID VARCHAR(50),
			Date DATE,
			Name TEXT,
			CVE TEXT,
			POC TEXT,
		    Description TEXT
			) DEFAULT CHARSET=utf8mb4;
		`)
	if err != nil {
		utils.PrintColor("error", "Error creating table packetstorm_db:", err, "\n")
		return
	}

}

func InsertPacketstormDB(id string, date time.Time, name string, cve string, poc string, desc string) error {
	_, err := packetstormdb.Exec(`
		INSERT INTO packetstorm_db (ID, Date, Name, CVE, POC, Description)
		VALUES (?, ?, ?, ?, ?, ?)
	`, id, date, name, cve, poc, desc)

	if err != nil {
		utils.PrintColor("error", "Error inserting vul %v into table packetstorm_db: %v\n", id, err)
		return err
	}

	return nil
}
