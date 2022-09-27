package database

import (
	"context"
	"therealbroker/internal/modules"
	"time"
)

func StoreToDatabase(ctx context.Context, row *modules.RowObject) {
	if db == nil {
		panic("NO DATABASE IS INITIALIZED!")
	}
	db.StoreToDatabase(ctx, row)
}

func GetMessage(id int) modules.RowObject {
	if db == nil {
		panic("NO DATABASE IS INITIALIZED!")
	}
	return db.GetMessage(id)
}


func InitPostgresDB() {
	(&postgresDB{}).InitDatabase()
}

func InitCassandraDB() {
	(&cassandraDB{}).InitDatabase()
}

// schedule runs function f once every t period of time
func schedule(f func(), t time.Duration) {
	for {
		f()
		time.Sleep(t)
	}
}
