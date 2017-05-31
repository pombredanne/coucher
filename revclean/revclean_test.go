package revclean

import (
	"testing"
	"net/http"
	"log"
	"github.com/yoavl/coucher/util"
	"github.com/yoavl/coucher/database"
	"fmt"
)

var db *database.Database

func init() {
	db = &database.Database{
		BaseUrl: "http://localhost:5984",
		Name:    "db1",
	}

	req, err := http.NewRequest("DELETE", db.Url(), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	str := util.ResponseToString(resp)
	fmt.Printf("%s", str)

	req, err = http.NewRequest("PUT", db.Url(), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	str = util.ResponseToString(resp)
	fmt.Printf("%s", str)
}

func TestConflicts(t *testing.T) {
	const docid = "1"

	docs := make([]map[string]interface{}, 0)

	doc1 := make(map[string]interface{})
	doc1["_id"] = docid
	doc1["msg"] = "wow"

	docs = append(docs, doc1)

	rev := writeDoc(docs)
	fmt.Println(rev)

	//Create conflicts by writing on top of the same rev with diff content
	docs = make([]map[string]interface{}, 0)

	doc2 := make(map[string]interface{})
	doc2["_id"] = docid
	doc2["_rev"] = rev
	doc2["msg"] = "bow"
	docs = append(docs, doc2)

	doc3 := make(map[string]interface{})
	doc3["_id"] = docid
	doc3["_rev"] = rev
	doc3["msg"] = "pow"
	docs = append(docs, doc3)

	doc4 := make(map[string]interface{})
	doc4["_id"] = docid
	doc4["_rev"] = rev
	doc4["msg"] = "low"
	docs = append(docs, doc4)

	rev = writeDoc(docs)
	fmt.Printf("First rev: %s\n", rev)

	resolvedRev, conflicts, err := db.GetRev(docid, true)
	if err != nil {
		log.Fatal(err)
	}
	if len(conflicts) != 2 {
		t.Error("Expected 2 conflicts! Resolved rev: %s. Conflicts: %v\n", resolvedRev, conflicts)
	}

	preservedRef, err := CleanRevs(db, docid, "", false)
	if err != nil {
		log.Fatal(err)
	}

	resolvedRev, conflicts, err = db.GetRev(docid, true)
	if err != nil {
		log.Fatal(err)
	}
	if len(conflicts) != 0 {
		t.Errorf("Expected 0 conflicts after cleanup! Resolved rev: %s. Conflicts: %v\n", resolvedRev, conflicts)
	}
	if preservedRef != resolvedRev {
		t.Errorf("Resolved rev %s is different from preserved rev %s\n", resolvedRev, preservedRef)
	}
}

func writeDoc(docs []map[string]interface{}) string {
	rev, err := db.WriteBulk(docs, true)
	if err != nil {
		log.Fatal(err)
	}
	return rev
}
