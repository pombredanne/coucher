package database

import (
	"net/http"
	"fmt"
	"bytes"
	"log"
	"github.com/yoav/coucher/util"
	"encoding/json"
	"github.com/bugsnag/bugsnag-go/errors"
	"strings"
)

const CT_JSON = "application/json"

const FLD_REQ_CONFLICTS = "conflicts"
const FLD_RES_CONFLICTS = "_conflicts"

type Database struct {
	BaseUrl  string `json:"url"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (db *Database) Url() string {
	if db.Username == "" {
		return fmt.Sprintf("%s/%s", db.BaseUrl, db.Name)
	} else {
		i := strings.Index(db.BaseUrl, ":")
		return fmt.Sprintf(
			"%s%s:%s@%s/%s", db.BaseUrl[:i+3], db.Username, db.Password, db.BaseUrl[i+3:], db.Name)

	}
}

func (db *Database) CleanConflicts(id string, keepRev string, dryrun bool) error {
	_, conflicts, err := db.GetRev(id, true)
	if err != nil {
		log.Fatal(err)
	}
	if conflicts == nil {
		log.Println("No conflicts for doc id: %s", id)
		return nil
	}
	for _, rev := range conflicts {
		url := fmt.Sprintf("%s/%s?rev=%s", db.Url(), id, rev)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		doc := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}
		doc["_deleted"] = true

		docs := make([]map[string]interface{}, 1)
		docs[0] = doc

		val, err := json.Marshal(docs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Updating: %s... ", val)
		if dryrun {
			fmt.Printf("Dry\n")
		} else {
			//Couch 1.0.2 does not support all_or_nothing :(
			db.WriteBulk(docs, false)
			fmt.Printf("Done\n")
		}
	}
	return nil
}

func (db *Database) GetRev(id string, withConflicts bool) (string, []string, error) {
	//curl "http://127.0.0.1:5984/db1/1?conflicts=true"
	url := fmt.Sprintf("%s/%s?%s=%t", db.Url(), id, FLD_REQ_CONFLICTS, withConflicts)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	str := util.ResponseToString(resp)
	//fmt.Printf("%s", str)

	doc := make(map[string]interface{})
	err = json.Unmarshal([]byte(str), &doc)
	rev, err := stringFromDoc(doc, "_rev")
	if err != nil {
		return "", nil, err
	}
	conflicts, err := conflictsFromDoc(doc)
	if err != nil {
		return "", nil, err
	}

	return rev, conflicts, nil
}

func (db *Database) WriteBulk(docs []map[string]interface{}, allOrNothing bool) (string, error) {
	url := fmt.Sprintf("%s/_bulk_docs", db.Url())

	req := make(map[string]interface{})
	req["all_or_nothing"] = allOrNothing
	req["docs"] = docs

	b, _ := json.Marshal(req)
	//fmt.Printf("%s", b)
	resp, err := http.Post(url, CT_JSON, bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	errMsg := util.UnsuccessfulMessage(resp)
	if errMsg != "" {
		log.Fatal(errMsg)
	}
	str := util.ResponseToString(resp)
	//fmt.Printf("%s", str)

	c := make([]map[string]interface{}, len(docs))
	err = json.Unmarshal([]byte(str), &c)
	if err != nil {
		log.Fatal(err)
	}
	return stringFromDoc(c[0], "rev")
}

func stringFromDoc(c map[string]interface{}, field string) (string, error) {
	//fmt.Printf("%v\n\n", c)
	raw := c[field]
	if raw == nil {
		return "", nil
	}
	if str, ok := raw.(string); ok {
		return str, nil
	} else {
		e := fmt.Sprintf("%v is not a string", raw)
		return "", errors.New(e, 1)
	}
}

func conflictsFromDoc(doc map[string]interface{}) ([]string, error) {
	raw := doc[FLD_RES_CONFLICTS]
	//fmt.Printf("CC: %v, type: %v\n\n", raw, reflect.TypeOf(raw))
	if raw == nil {
		return make([]string, 0), nil
	}
	if array, ok := raw.([]interface{}); ok {
		ret := make([]string, len(array))
		for i, v := range array {
			if str, ok := v.(string); ok {
				ret[i] = str
			} else {
				e := fmt.Sprintf("%v is not a string", raw)
				return nil, errors.New(e, 1)
			}

		}
		return ret, nil
	} else {
		e := fmt.Sprintf("%v is not an array", raw)
		return nil, errors.New(e, 1)
	}
}
