package main

import (
	"testing"
	"os"
)

func TestRevClean(t *testing.T) {
	os.Setenv("COUCH_URL", "http://localhost:5984")
	os.Setenv("COUCH_DBNAME", "db1")
	os.Setenv("COUCH_USER", "")
	os.Setenv("COUCH_PASS", "")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"coucher", "revclean","--docid=1", "--dry=true"}
	main()
}
