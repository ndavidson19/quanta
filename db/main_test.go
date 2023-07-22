package db

import ("database/sql"
	"testing"
	"os"
	"log"

	_ "github.com/lib/pq")

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/go_client?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: %w", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}