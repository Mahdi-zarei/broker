package database

import (
	"context"
	"therealbroker/internal/modules"
)

var db Database

type Database interface {
	ConnectToDb()

	StoreToDatabase(ctx context.Context, row *modules.RowObject)

	GetMessage(id int) modules.RowObject

	InitDatabase()
}

func setDatabase(database Database) {
	db = database
}
