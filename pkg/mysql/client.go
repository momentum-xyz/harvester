package mysql

import (
	"database/sql"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func NewDB(cfg *Config) (*sql.DB, func(), error) {
	mysqlConfig := makeMYSQLConfig(cfg)
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup, nil
}

func GetMYSQLConfig(cfg *Config) *mysql.Config {
	return makeMYSQLConfig(cfg)
}

func makeMYSQLConfig(cfg *Config) *mysql.Config {
	mysqlConfig := mysql.NewConfig()

	mysqlConfig.User = cfg.Username
	mysqlConfig.Passwd = cfg.Password
	mysqlConfig.DBName = cfg.Database
	mysqlConfig.Addr = cfg.Host + ":" + strconv.FormatUint(uint64(cfg.Port), 10)
	mysqlConfig.Net = "tcp"
	mysqlConfig.InterpolateParams = true
	mysqlConfig.ParseTime = true

	return mysqlConfig
}
