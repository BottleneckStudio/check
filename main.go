package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BottleneckStudio/check/check"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "check"
	app.Usage = "check if a website is up or down!"
	app.Action = func(c *cli.Context) error {
		site := c.Args().Get(0)
		isUp := check.IsUp(site)
		var msg string
		if isUp {
			msg = "up"
		} else {
			msg = "down"
		}
		fmt.Printf("Site %s is %s", site, msg)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
