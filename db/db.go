package db

import (
	"fmt"
	"packagelock/config"

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
		errorMessage := fmt.Sprintf(` Couldn't connect to DB! Got: '%s'.
		1. Check the config for a wrong Address/Port (Currently: %s:%s)
		2. Check if the DB is reachable (eg. a Ping). Check the Firewalls if there.
		3. Consult the PackageLock Doc's! 🚀
		Golang Trace Logs:
			`, err.Error(), dbAddress, dbPort)
		fmt.Println(errorMessage)
		panic(err)
	}

	if _, err = db.Signin(map[string]interface{}{
		// TODO: get user&password from Conf
		"user": dbUsername,
		"pass": dbPasswd,
	}); err != nil {
		// FIXME: Logging of wrong username and maybe SHA-PASSWD?

		errorMessage := fmt.Sprintf(` Couldn't connect to DB! Got: '%s'.
		1. Check the config for a wrong DB-Username/Password (Currently: %s/<read the config!>)
		3. Consult the PackageLock Doc's! 🚀
		Golang Trace Logs:
			`, err.Error(), dbUsername)
		fmt.Println(errorMessage)
		panic(err)
	}

	if _, err = db.Use("PackageLock", "db1.0"); err != nil {
		// No error handling possible, as we need to use this db
		panic(err)
	}

	DB = db

	return nil
}

// If you use this, fix it!
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
