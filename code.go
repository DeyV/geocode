package main

import (
	"encoding/json"
	"fmt"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"sync"
)

var cache = make(map[string]*CodeResponse)
var lock sync.Mutex

type CodeResponse struct {
	Kod    string
	Gmina  string
	Powiat string

	Wojewodztwo   string
	WojewodztwoId int
}

func codeJsonHandler(rw http.ResponseWriter, req *http.Request) {
	// fmt.Print("+")
	code := req.FormValue("code")

	r := getCode(code)
	m, err := json.Marshal(r)

	if err != nil {
		fmt.Fprintln(rw, "error:", err)
	}

	rw.Write(m)
}

func codeTextHandler(rw http.ResponseWriter, req *http.Request) {
	// fmt.Print("+")
	code := req.FormValue("code")
	r := getCode(code)

	fmt.Fprintf(rw, "%s;%s", r.Kod, r.Powiat)
}

func codeTextCacheHandler(rw http.ResponseWriter, req *http.Request) {
	// fmt.Print("+")
	code := req.FormValue("code")

	lock.Lock()
	defer lock.Unlock()

	r, ok := cache[code]
	if !ok {
		r = getCode(code)
		cache[code] = r
	}

	fmt.Fprintf(rw, "%s;%s", r.Kod, r.Powiat)
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "resources/teryt.sqlite")
	if err != nil {
		panic(err)
	}
}

func getCode(code string) *CodeResponse {
	rows, err := db.Query(fmt.Sprintf("SELECT kod, powiat, gmina, wojewodztwo FROM poczta where kod = '%s'", code))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer rows.Close()

	cr := &CodeResponse{}
	for rows.Next() {
		rows.Scan(&cr.Kod, &cr.Powiat, &cr.Gmina, &cr.WojewodztwoId)
		return cr
	}

	return nil
}
