package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/codegangsta/cli"

	"github.com/yuuki1/grabeni/commands"
	"github.com/yuuki1/grabeni/log"
)

var AppHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

{{.Usage}}

Version: {{.Version}}{{if or .Author .Email}}

Author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

var commandArgs = map[string]string{
	"status": commands.CommandArgStatus,
	"list":   commands.CommandArgList,
	"attach": commands.CommandArgAttach,
	"detach": commands.CommandArgDetach,
	"grab":   commands.CommandArgGrab,
}

func setDebugOutputLevel() {
	for _, f := range os.Args {
		if f == "-D" || f == "--debug" || f == "-debug" {
			log.IsDebug = true
		}
	}

	debugEnv := os.Getenv("GRABENI_DEBUG")
	if debugEnv != "" {
		showDebug, err := strconv.ParseBool(debugEnv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing boolean value from GRABENI_DEBUG: %s\n", err)
			os.Exit(1)
		}
		log.IsDebug = showDebug
	}
}

func init() {
	setDebugOutputLevel()
	argsTemplate := "{{if false}}"
	for _, command := range append(commands.Commands) {
		argsTemplate = argsTemplate + fmt.Sprintf("{{else if (eq .Name %q)}}%s %s", command.Name, command.Name, commandArgs[command.Name])
	}
	argsTemplate = argsTemplate + "{{end}}"

	cli.CommandHelpTemplate = `Usage: grabeni ` + argsTemplate + `

{{.Usage}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .Flags}}

Options:
   {{range .Flags}}
   {{.}}{{end}}{{ end }}
`

	cli.AppHelpTemplate = AppHelpTemplate
}

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "y_uuki"
	app.Email = "https://github.com/yuuki1/grabeni"
	app.Commands = commands.Commands
	app.CommandNotFound = cmdNotFound
	app.Usage = "An ops-friendly AWS ENI grabbing tool"
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Enable debug mode",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Logf("error", err.Error())
	}
}

func cmdNotFound(c *cli.Context, command string) {
	log.Logf(
		"",
		"%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name,
		command,
		c.App.Name,
		os.Args[0],
	)
}
