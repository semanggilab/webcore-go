package mongo

import (
	"github.com/semanggilab/webcore-go/app/loader"
)

type MongoLoader struct {
}

func (a *MongoLoader) ClassName() string {
	return "MongoDatabase"
}

func (l *MongoLoader) Init(args ...any) (loader.Library, error) {
	// config := args[1].(config.DatabaseConfig)

	db := &MongoDatabase{}
	err := db.Install(args...)
	if err != nil {
		return nil, err
	}

	err = db.Connect()
	if err != nil {
		return nil, err
	}

	return db, nil
}
