package util

import (
	"io/ioutil"
	"net/http"
	"log"
	"fmt"
)

func ResponseToString(resp *http.Response) string {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body[:])
}

func UnsuccessfulMessage(resp *http.Response) string {
	if resp.StatusCode >= 400 {
		return fmt.Sprintf("%d: %s", resp.StatusCode, ResponseToString(resp))
	} else {
		return ""
	}
}
