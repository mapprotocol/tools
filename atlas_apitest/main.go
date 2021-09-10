package main

import (
	"fmt"
	"github.com/mapprotocol/tools/ethclient"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var (
	// The app that holds all commands and flags.
	app *cli.App
)

func init() {
	app = cli.NewApp()
	app.Usage = "Atlas  Tool"
	app.Name = filepath.Base(os.Args[0])
	app.Version = "1.0.0"
	app.Copyright = "Copyright 2020-2021 The Atlas Authors"
	app.Flags = []cli.Flag{}
	app.Action = MigrateFlags(start)
	app.CommandNotFound = func(ctx *cli.Context, cmd string) {
		fmt.Fprintf(os.Stderr, "No such command: %s\n", cmd)
		os.Exit(1)
	}
	// Add subcommands.
	app.Commands = []cli.Command{}
	cli.CommandHelpTemplate = OriginCommandHelpTemplate
	sort.Sort(cli.CommandsByName(app.Commands))
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var OriginCommandHelpTemplate string = `{{.Name}}{{if .Subcommands}} command{{end}}{{if .Flags}} [command options]{{end}} [arguments...] {{if .Description}}{{.Description}} {{end}}{{if .Subcommands}} SUBCOMMANDS:     {{range .Subcommands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}     {{end}}{{end}}{{if .Flags}} OPTIONS: {{range $.Flags}}{{"\t"}}{{.}} {{end}} {{end}}`

func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.GlobalSet(name, ctx.String(name))
			}
		}
		return action(ctx)
	}
}
func start(ctx *cli.Context) error {
	BlockNumber()
	select {}
	return nil
}
func dialConn() (*ethclient.Client, string) {
	ip := "119.8.165.158"
	port := 8082
	url := fmt.Sprintf("http://%s", fmt.Sprintf("%s:%d", ip, port))
	conn, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the atlas chain client: %v", err)
	}
	return conn, url
}
