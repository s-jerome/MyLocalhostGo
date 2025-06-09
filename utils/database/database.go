package utils

import (
	"database/sql"
	"mylocalhost/utils"
)

func OpenSQLiteConnection(dbFilePath string) (*sql.DB, bool, error) {
	var dbFileExists, fileExistsErr = utils.FileExists(dbFilePath)
	if fileExistsErr != nil {
		return nil, false, fileExistsErr
	}

	var db, openErr = sql.Open("sqlite3", dbFilePath)
	return db, dbFileExists, openErr
}
