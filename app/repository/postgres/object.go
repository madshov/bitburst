package postgres

import (
	"database/sql"
	"strconv"

	"github.com/madshov/bitburst/app"
)

type ObjectRepo struct {
	DB *sql.DB
}

func NewObjectRepo(DB *sql.DB) *ObjectRepo {
	return &ObjectRepo{
		DB: DB,
	}
}

// CreateObjects creates objects if the don't exist, otherwise the timestamp
// is updated.
func (or *ObjectRepo) CreateObjects(objs []app.Object) error {
	sqlStr := "INSERT INTO objects(object_id, online) VALUES ($1, $2) ON CONFLICT (object_id) DO UPDATE SET online = $2, timestamp = current_timestamp"

	for _, obj := range objs {
		//prepare the statement
		stmt, err := or.DB.Prepare(sqlStr)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(obj.ID, strconv.FormatBool(obj.Online))
		if err != nil {
			return err
		}

		defer stmt.Close()
	}

	return nil
}

// DeleteObjects deletes objects older than 30 seconds.
func (or *ObjectRepo) DeleteObjects() error {
	sqlStr := "DELETE FROM objects WHERE timestamp  < now() - '30 seconds' :: interval"

	_, err := or.DB.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}
