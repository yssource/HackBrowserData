package main

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	version           string
	browserName       string
	exportDir         string
	outputFormat      string
	verbose           bool
	compress          bool
	customProfilePath string
	customKeyPath     string
)

func main() {
}

func Execute() {
	app := &cli.App{
		Name:  "hack-browser-data",
		Usage: "Export passwords/cookies/history/bookmarks from browser",
		UsageText: "[hack-browser-data -b chrome -f json -dir results -cc]\n 	Get all data(password/cookie/history/bookmark) from chrome",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"vv"}, Destination: &verbose, Value: false, Usage: "verbose"},
			&cli.BoolFlag{Name: "compress", Aliases: []string{"cc"}, Destination: &compress, Value: false, Usage: "compress result to zip"},
			&cli.StringFlag{Name: "browser", Aliases: []string{"b"}, Destination: &browserName, Value: "all", Usage: "available browsers: all|" + strings.Join(core.ListBrowser(), "|")},
			&cli.StringFlag{Name: "results-dir", Aliases: []string{"dir"}, Destination: &exportDir, Value: "results", Usage: "export dir"},
			&cli.StringFlag{Name: "format", Aliases: []string{"f"}, Destination: &outputFormat, Value: "csv", Usage: "format, csv|json|console"},
			&cli.StringFlag{Name: "profile-dir-path", Aliases: []string{"p"}, Destination: &customProfilePath, Value: "", Usage: "custom profile dir path, get with chrome://version"},
			&cli.StringFlag{Name: "key-file-path", Aliases: []string{"k"}, Destination: &customKeyPath, Value: "", Usage: "custom key file path"},
		},
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
