package db

import "github.com/surrealdb/surrealdb.go"

var DB *surrealdb.DB

func InitDB() error {
	db, err := surrealdb.New("ws://localhost:8000/rpc")
	if err != nil {
		panic(err)
	}

	if _, err = db.Signin(map[string]interface{}{
		// TODO: get user&password from Conf
		"user": "root",
		"pass": "root",
	}); err != nil {
		panic(err)
	}

	if _, err = db.Use("PackageLock", "db1.0"); err != nil {
		// This Error indecates the non existance of ther the
		// PackageLock Namespace or the db.
		// If this happens, we should run a basic db setup.
		// TODO: Create DB migration in migration.go
		// Ether way we should log and handle the error
		// FIXME: Logging
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
