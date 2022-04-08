package main

import (
	"log"

	"github.com/ClickerMonkey/sudogo/pkg/rest"
)

func main() {
	log.Fatal(rest.Start(":3000"))
}
