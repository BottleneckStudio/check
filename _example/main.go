package main

import (
	"log"

	"github.com/BottleneckStudio/check/check"
)

func main() {

	isUp := check.IsUp("google.com")

	log.Println(isUp)
}
