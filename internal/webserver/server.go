package webserver

import (
	"log"
	"net/http"
)

// Start starts the web server.
func Start() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Println("Serving web UI on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
