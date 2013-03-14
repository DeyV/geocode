package main

import (
	"encoding/json"
	"fmt"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"sync"
)

func codeHandler(rw http.ResponseWriter, req *http.Request) {
	// fmt.Print("+")
	code := req.FormValue("code")

	k := &CodeResponse{
		Kod:   getCode(code),
		Gmina: code,
	}

	m, err := json.Marshal(k)

	if err != nil {
		fmt.Fprintln(rw, "error:", err)
	}

	rw.Write(m)
}

type CodeResponse struct {
	Kod    string
	Gmina  string
	Powiat string

	Wojewodztwo   string
	WojewodztwoId int
}

func codeTextHandler(rw http.ResponseWriter, req *http.Request) {
	// fmt.Print("+")
	code := req.FormValue("code")

	fmt.Fprintf(rw, "%s;%s", getCode(code), code)
}

var cache = make(map[string]string)
var lock sync.Mutex

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

	fmt.Fprintf(rw, "%s;%s", r, code)
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "resources/teryt.sqlite")
	if err != nil {
		panic(err)
	}
}

func getCode(code string) string {
	rows, err := db.Query(fmt.Sprintf("SELECT id, kod FROM poczta where kod = '%s'", code))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var kod string
		rows.Scan(&id, &kod)
		return kod
	}

	return ""
}
