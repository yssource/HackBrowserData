package main

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	hbd "github.com/moond4rk/hack-browser-data"
)

var (
	browserName  string
	itemName     string
	exportDir    string
	outputFormat string
	verbose      bool
	compress     bool
)

func main() {
	app := &cli.App{
		Name:  "hack-browser-data",
		Usage: "Export passwords/cookies/history/bookmarks from browser",
		UsageText: "[hack-browser-data -b chrome -f json -dir results -cc]\n 	Get all data(password/cookie/history/bookmark) from chrome",
		Version: "v0.4.0",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"vv"}, Destination: &verbose, Value: false, Usage: "verbose"},
			&cli.BoolFlag{Name: "compress", Aliases: []string{"cc"}, Destination: &compress, Value: false, Usage: "compress result to zip"},
			&cli.StringFlag{Name: "browser", Aliases: []string{"b"}, Destination: &browserName, Value: "all", Usage: "available browsers: all|" + strings.Join(hbd.ListAllBrowserName(), "|")},
			&cli.StringFlag{Name: "item", Aliases: []string{"i"}, Destination: &itemName, Value: "all", Usage: "available item: all|" + strings.Join(hbd.ListAllItemerName(), "|")},
			&cli.StringFlag{Name: "output-dir", Aliases: []string{"dir"}, Destination: &exportDir, Value: "results", Usage: "export dir"},
			&cli.StringFlag{Name: "format", Aliases: []string{"f"}, Destination: &outputFormat, Value: "csv", Usage: "format, csv|json"},
		},
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			var options []hbd.Option
			if browserName == "all" {
				options = append(options, hbd.WithAllClients())
			} else {
				options = append(options, hbd.WithClient(browserName))
			}
			if itemName == "all" {
				options = append(options, hbd.WithAllItemers())
			} else {
				options = append(options, hbd.WithItemer(itemName))
			}

			if outputFormat == "json" {
				options = append(options, hbd.WithOutputJson())
			} else {
				options = append(options, hbd.WithOutputCSV())
			}

			options = append(options, hbd.WithOutputDir(exportDir))
			options = append(options, hbd.WithCompress(compress))
			browser, err := hbd.NewBrowser(options...)
			if err != nil {
				panic(err)
			}
			if err := browser.Run(); err != nil {
				panic(err)
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
