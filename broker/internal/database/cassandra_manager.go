package database

import (
	"context"
	"github.com/gocql/gocql"
	"go.opentelemetry.io/otel"
	"log"
	"therealbroker/internal/modules"
	"time"
)

type cassandraDB struct {
	session *gocql.Session
}

func (cs *cassandraDB) StoreToDatabase(ctx context.Context, row *modules.RowObject) {
	_, span := otel.Tracer(name).Start(ctx, "dbStore")
	defer span.End()

	err := cs.session.Query("INSERT INTO messages (identifier, body, creation, expiration, subject) VALUES (?, ?, ?, ?, ?)",
		row.Id, row.Body, row.Creation, row.Expiration, row.Subject).Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func (cs *cassandraDB) ConnectToDb() {
	cluster := gocql.NewCluster("cassandra")
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second

	var err error
	tmp, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	err = tmp.Query("CREATE KEYSPACE IF NOT EXISTS broker WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1}").Exec()
	if err != nil {
		log.Fatal(err)
	}
	tmp.Close()

	cluster.Keyspace = "broker"
	cs.session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	err = cs.session.Query("CREATE TABLE IF NOT EXISTS messages(identifier int, body text, creation timestamp," +
		" expiration int, subject text, PRIMARY KEY ( (subject), identifier) )").Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func (cs *cassandraDB) GetMessage(id int) modules.RowObject {
	var (
		body string
		cr   time.Time
		exp  int
		subj string
	)
	err := cs.session.Query("SELECT body, creation, expiration, subject FROM messages WHERE identifier = ?", id).
		Consistency(gocql.One).Scan(&body, &cr, &exp, &subj)
	if err != nil {
		return *modules.NewRowObject(-1, "", time.Now(), -1, "")
	}
	return *modules.NewRowObject(id, body, cr, exp, subj)
}

func (cs *cassandraDB) InitDatabase() {
	cs.ConnectToDb()
	setDatabase(cs)
}
