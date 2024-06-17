package main

import (
	"log"
	"os"

	"github.com/common-fate/linuxpack/cmd/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "linuxpack",
		Usage: "Package and publish an APT repository to S3 and CloudFront",
		Commands: []*cli.Command{
			&command.Package,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
