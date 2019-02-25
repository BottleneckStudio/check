package main

import (
	"log"

	"github.com/check"
)

func main() {

	isUp := check.IsUp("google.com")

	log.Println(isUp)
}
