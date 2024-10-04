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

	if _, err = db.Use("test", "test"); err != nil {
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
	return nil
}
