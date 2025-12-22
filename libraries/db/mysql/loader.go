package mysql

import (
	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/loader"
	"github.com/semanggilab/webcore-go/lib/sql"
	"gorm.io/driver/mysql"
)

type MysqlLoader struct {
	// DB *sql.SQLDatabase
}

func (a *MysqlLoader) ClassName() string {
	return "SQLDatabase"
}

func (l *MysqlLoader) Init(args ...any) (loader.Library, error) {
	config := args[1].(config.DatabaseConfig)
	dsn := sql.BuildDSN(config)

	db := &sql.SQLDatabase{}
	db.SetDialect(mysql.Open(dsn))
	err := db.Install(args...)
	if err != nil {
		return nil, err
	}

	db.Connect()

	// l.DB = db
	return db, nil
}
