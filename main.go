package main

import (
	"log"
	"net/http"

	"github.com/steffbo/intercise/router"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", router.Router()))
}
