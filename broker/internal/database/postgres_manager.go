package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"log"
	"strings"
	"sync"
	"therealbroker/Internal/modules"
	"time"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "broker"
	name     = "postgres"
)

type postgresDB struct {
	session *sql.DB

	strs      []string
	args      []interface{}
	waitermng modules.WaiterManager
	ready     bool

	syncer sync.Mutex
}

func (psg *postgresDB) ConnectToDb() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database: ", psqlInfo)
	psg.session = db
}

func (psg *postgresDB) executeBatch() {
	if !psg.ready {
		return
	}
	block, span := otel.Tracer(name).Start(context.Background(), "commit")
	defer span.End()

	_, span2 := otel.Tracer(name).Start(block, "block time")
	psg.syncer.Lock()
	span2.End()

	stmt := "INSERT INTO messages (identifier, body, creation, expiration, subject) VALUES " + strings.Join(psg.strs, ",")

	res, err := psg.session.Exec(stmt, psg.args...)
	if err != nil {
		log.Fatal(res, err)
	}
	log.Print("Committed changes to Database size: ", len(psg.args)/5)

	psg.strs = make([]string, 0)
	psg.args = make([]interface{}, 0)

	psg.waitermng.FreeWaiters()

	psg.ready = false

	psg.syncer.Unlock()

}

func (psg *postgresDB) waitForExec(row *modules.RowObject) {
	psg.syncer.Lock()

	sz := len(psg.strs)

	psg.strs = append(psg.strs, fmt.Sprintf("($%v, $%v, $%v, $%v, $%v)", sz*5+1, sz*5+2, sz*5+3, sz*5+4, sz*5+5))

	psg.args = append(psg.args, row.Id, row.Body, row.Creation, row.Expiration, row.Subject)

	tmpChan := psg.waitermng.RegisterWaiter()

	psg.ready = true

	psg.syncer.Unlock()

	<-tmpChan
}

func (psg *postgresDB) StoreToDatabase(ctx context.Context, row *modules.RowObject) {
	_, span := otel.Tracer(name).Start(ctx, "RecordToDB")
	defer span.End()

	psg.waitForExec(row)
}

func (psg *postgresDB) GetMessage(id int) modules.RowObject {
	rows, err := psg.session.Query("SELECT body, creation, expiration, subject FROM messages WHERE identifier = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	var (
		body     string
		creation time.Time
		exp      int
		subj     string
	)
	for rows.Next() {
		rows.Scan(&body, &creation, &exp, &subj)
		return *modules.NewRowObject(-1, "", time.Now(), -1, "")
	}
	return *modules.NewRowObject(id, body, creation, exp, subj)
}

func (psg *postgresDB) InitDatabase() {
	psg.ConnectToDb()
	psg.syncer = sync.Mutex{}
	psg.ready = false
	psg.waitermng = *modules.NewWaiterManager()

	setDatabase(psg)

	go schedule(psg.executeBatch, 4*time.Millisecond)
}
