package database

import (
	"crush/config"
	"crush/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var seebugdb *sql.DB

func init() {
	conf := config.LoadConfig()
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s", conf.Database.DBUsername, conf.Database.DBPassword, conf.Database.DBPort, conf.Database.Name)
	var err error
	seebugdb, err = sql.Open("mysql", connectionString)
	if err != nil {
		utils.PrintColor("error", "Error connecting mysql database:", err, "\n")
		return
	}
	seebugdb.SetMaxIdleConns(20)
	seebugdb.SetMaxOpenConns(50)

	if err := seebugdb.Ping(); err != nil {
		utils.PrintColor("error", "Error connecting mysql database:", err, "\n")
		return
	}
}

func CreateSeebugDB() {
	_, err := seebugdb.Exec(`
		CREATE TABLE IF NOT EXISTS seebug_db (
			SSVID VARCHAR(50),
			SubmitTime DATE,
		    Severity VARCHAR(50),
			Name TEXT,
			CVE TEXT,
			HasPOC BOOLEAN
		) DEFAULT CHARSET=utf8mb4;
	`)
	if err != nil {
		utils.PrintColor("error", "Error creating table seebug_db:", err, "\n")
		return
	}

}

func InsertSeebugDB(ssvid string, submittime time.Time, severity string, name string, cve string, haspoc bool) error {
	_, err := seebugdb.Exec(`
		INSERT INTO seebug_db (SSVID, SubmitTime, Severity, Name, CVE, HasPoc)
		VALUES (?, ?, ?, ?, ?, ?)
	`, ssvid, submittime, severity, name, cve, haspoc)

	if err != nil {
		utils.PrintColor("error", "Error inserting vul %v into table seebug_db: %v\n", ssvid, err)
		return err
	}

	return nil
}
