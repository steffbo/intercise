package main

import (
	"net/http"

	"github.com/steffbo/intercise/router"
)

func main() {
	router.Router()
	http.ListenAndServe(":8080", nil)
}
