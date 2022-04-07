package main

import (
	"log"
	"net/http"

	"github.com/ClickerMonkey/sudogo/pkg/rest"
)

func main() {
	rest.Register()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
