package psql

import (
	"database/sql"
	"io/ioutil"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}