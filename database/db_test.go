package database

import (
	"testing"
)

func TestDb(t *testing.T) {
	db := &Database{
		BaseUrl: "http://localhost:5984",
		Name:    "db1",
	}

	if (db.Url() != "http://localhost:5984/db1") {
		t.Errorf("Wrong URL: %s\n", db.Url())
	}

	db = &Database{
		BaseUrl: "http://localhost:5984",
		Name:    "db1",
		Username: "shimi",
	}

	if (db.Url() != "http://shimi:@localhost:5984/db1") {
		t.Errorf("Wrong URL: %s\n", db.Url())
	}

	db = &Database{
		BaseUrl: "https://localhost:5984",
		Name:    "db1",
		Username: "shimi",
		Password: "secR3t",
	}

	if (db.Url() != "https://shimi:secR3t@localhost:5984/db1") {
		t.Errorf("Wrong URL: %s\n", db.Url())
	}
}
