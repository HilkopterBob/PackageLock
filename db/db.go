package db

import (
	"packagelock/config"
	"packagelock/logger"

	"github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

func InitDB() error {
	dbAddress := config.Config.GetString("database.address")
	dbPort := config.Config.GetString("database.port")
	dbUsername := config.Config.GetString("database.username")
	dbPasswd := config.Config.GetString("database.password")

	db, err := surrealdb.New("ws://" + dbAddress + ":" + dbPort + "/rpc")
	if err != nil {
		logger.Logger.Errorf(` Couldn't connect to DB! Got: '%s'.
		1. Check the config for a wrong Address/Port (Currently: %s:%s)
		2. Check if the DB is reachable (eg. a Ping). Check the Firewalls if there.
		3. Consult the PackageLock Doc's! 🚀
		Golang Trace Logs:
			`, err.Error(), dbAddress, dbPort)
	}

	if _, err = db.Signin(map[string]interface{}{
		"user": dbUsername,
		"pass": dbPasswd,
	}); err != nil {
		logger.Logger.Errorf(` Couldn't connect to DB! Got: '%s'.
		1. Check the config for a wrong DB-Username/Password (Currently: %s/<read the config!>)
		3. Consult the PackageLock Doc's! 🚀
		Golang Trace Logs:
			`, err.Error(), dbUsername)
	}

	if _, err = db.Use("PackageLock", "db1.0"); err != nil {
		// No error handling possible, as we need to use this db
		logger.Logger.Panicf("Couldn't Use 'PackageLock' Namespace and 'db1.0' Database. Got: %s", err)
	}

	DB = db

	logger.Logger.Infof("Successfully Connected to DB, at: %s:%s", dbAddress, dbPort)
	return nil
}

// INFO:  If you use this, fix it!
// INFO: 	And add logging/error handling
func Select(tablename string, SliceOfType interface{}) error {
	transaction, err := DB.Select(tablename)
	if err != nil {
		// FIXME: logging?
		// Error handling
		panic(err)
	}

	err = surrealdb.Unmarshal(transaction, &SliceOfType)
	if err != nil {
		// FIXME: Logging?
		// Error Handling?
		panic(err)
	}

	// FIXME: Add Success msg in Log!
	return nil
}
