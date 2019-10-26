package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/thesephist/hurricane/pkg/airtable"
)

func main() {
	base := airtable.Base{
		Id:     os.Getenv("HURRICANE_BASE_ID"),
		ApiKey: os.Getenv("HURRICANE_API_KEY"),
	}

	table := airtable.Table{
		Name: "Hackathon Applications",
		Base: base,
	}

	cache := airtable.Cache{Limit: 60 * time.Second} // cache responses for 1m
	cache.Purge()

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header.Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, cache.Get("/", func() string {
			return table.View("Approved")
		}))
	})
	r.HandleFunc("/{recordId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Header.Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, cache.Get(r.RequestURI, func() string {
			return table.Get(vars["recordId"])
		}))
	})
	r.HandleFunc("/view/{viewName}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Header.Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, cache.Get(r.RequestURI, func() string {
			return table.View(vars["viewName"])
		}))
	})
	http.Handle("/", r)

	fmt.Printf("Starting Hurricane proxy for %s\n", base.Id)

	port := os.Getenv("HURRICANE_PORT")
	if port == "" {
		port = "6070"
	}

	http.ListenAndServe(":"+port, nil)
}
