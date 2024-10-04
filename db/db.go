package db

import (
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
		panic(err)
	}

	if _, err = db.Signin(map[string]interface{}{
		// TODO: get user&password from Conf
		"user": dbUsername,
		"pass": dbPasswd,
	}); err != nil {
		// FIXME: Error handling
		// FIXME: Logging of wrong username and maybe SHA-PASSWD?
		panic(err)
	}

	if _, err = db.Use("PackageLock", "db2.0"); err != nil {
		// This Error indecates the non existance of ther the
		// PackageLock Namespace or the db.
		// If this happens, we should run a basic db setup.
		// TODO: Create DB migration in migration.go
		// Ether way we should log and handle the error
		// FIXME: Error handling
		panic(err)
	}

	DB = db

	return nil
}

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
