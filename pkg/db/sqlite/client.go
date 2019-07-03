package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github_statistics/internal/config"
	"github_statistics/internal/log"
	"os"
	"os/exec"
	"time"
)

var db *sql.DB
var dbName = fmt.Sprintf("%s/source.db", config.NewServerConfig().Http.DataDir)

func InitSqliteClient(createCallback func()) {
	log.Infof("preparing sqlite database:%s", dbName)
	_, err := os.Open(dbName)
	if err != nil {
		log.Warningf("database:%s does not exist and is about to create", dbName)
		_, err = os.Create(dbName)
		if err != nil {
			log.Fatalf("create database:%s found error", dbName, err)
		}
	}
	db, err = sql.Open("sqlite3", dbName)
	createCallback()
	if err != nil {
		log.Fatalf("init sqlite database:%s found error:%v", dbName, err)
	}
}

func CreateTable(table string) {
	log.Infof("[SQLITE SQL]%s", table)
	stmt, _ := db.Prepare(table)
	_, err := stmt.Exec()
	if err != nil {
		log.Fatalf("exec statement:%s found error:%v", table, err)
	}
}

func Query(query string) error {
	log.Infof("[EXEC SQL]%s", query)
	stmt, _ := db.Prepare(query)
	_, err := stmt.Exec()
	if err != nil {
		log.Errorf("insert data found error:%v", err)
		return err
	}
	return nil
}

func genBakDbName(year, month int) (string, error) {
	dataDir := config.NewServerConfig().Http.DataDir
	bakDbName := fmt.Sprintf("%s/%d-%02d.bak.db", dataDir, year, month)
	_, err := os.Create(bakDbName)
	if err != nil {
		log.Errorf("create backup database:%s found error:%v", bakDbName, err)
		return "", nil
	}
	return bakDbName, nil
}

var tables = []string{"developer"}

func StartBak() {
	log.Infof("start to backup data")
	for _, table := range tables {
		log.Infof("backup database:%s", table)
		err := bakTable(table)
		if err != nil {
			log.Errorf("backup table:%s found error:%v", table, err)
		}
	}
	log.Info("database backup finished")

}

func bakTable(table string) error {
	type createdAt struct {
		CreatedAt time.Time
	}
	var rets []createdAt
	var first time.Time
	var last time.Time

	rows, _ := db.Query(fmt.Sprintf("SELECT created_at FROM %s ORDER BY created_at DESC LIMIT 1", table))
	err := rows.Scan(&rets)
	if err != nil {
		return err
	}
	if len(rets) == 0 {
		return nil
	}
	first = rets[0].CreatedAt

	rows, _ = db.Query(fmt.Sprintf("SELECT created_at FROM %s ORDER BY created_at ASC LIMIT 1", table))
	if err != nil {
		return err
	}
	last = rets[0].CreatedAt

	if first.Year() == last.Year() &&
		(int(first.Month())-int(last.Month()) < 1) {
		return nil
	}

	err = doBak(&first, &last)
	if err != nil {
		return err
	}
	return delBakRecords(table, &first)
}

func delBakRecords(table string, first *time.Time) error {
	year := first.Year()
	month := int(first.Month())

	condi := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
	_, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE created_at < '%s'", table, condi))
	return err
}

func doBak(first, last *time.Time) error {
	year := first.Year()
	month := int(first.Month())
	for {
		if year == last.Year() && month < int(last.Month()) {
			break
		}
		err := dumpData(year, month)
		if err != nil {
			return err
		}

		month -= 1
		if month < 1 {
			year -= 1
			month = 12
		}
	}
	return nil
}

func dumpData(year, month int) error {
	bakDbName, err := genBakDbName(year, month)
	if err != nil {
		return err
	}
	outs, err := exec.Command("/bin/bash", "-c", fmt.Sprint("%s .dump %s", dbName, bakDbName)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(outs))
	}
	return nil
}
