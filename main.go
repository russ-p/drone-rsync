package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	_ "github.com/joho/godotenv/autoload"
)

var (
	buildCommit string
	version     string // build number set at compile-time
)

func main() {
	fmt.Printf("Drone Rsync Plugin built from %s\n", buildCommit)

	app := cli.NewApp()
	app.Name = "Rsync"
	app.Usage = "Uses rsync to upload files to a remote server"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "host",
			Usage:  "connect to host",
			EnvVar: "PLUGIN_HOST,SSH_HOST",
		},
		cli.StringFlag{
			Name:   "user",
			Usage:  "connect as user",
			EnvVar: "PLUGIN_USER,SSH_USER",
			Value:  "root",
		},
		cli.IntFlag{
			Name:   "port",
			Usage:  "connect to port",
			EnvVar: "PLUGIN_PORT,SSH_PORT",
			Value:  22,
		},
		cli.StringFlag{
			Name:   "source",
			Usage:  "source path from which files are copied",
			EnvVar: "PLUGIN_SOURCE,SOURCE",
			Value:  "./",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "target path to which files are copied",
			EnvVar: "PLUGIN_TARGET,TARGET",
			Value:  "/",
		},
		cli.BoolFlag{
			Name:   "delete",
			Usage:  "delete extraneous files from the target dir",
			EnvVar: "PLUGIN_DELETE,DELETE",
		},
		cli.BoolFlag{
			Name:   "recursive",
			Usage:  "recursively transfer all files",
			EnvVar: "PLUGIN_RECURSIVE,RECURSIVE",
		},
		cli.StringSliceFlag{
			Name:   "include",
			Usage:  "include files matching the specified pattern",
			EnvVar: "PLUGIN_INCLUDE,INCLUDE",
		},
		cli.StringSliceFlag{
			Name:   "exclude",
			Usage:  "exclude files matching the specified pattern",
			EnvVar: "PLUGIN_EXCLUDE,EXCLUDE",
		},
		cli.StringSliceFlag{
			Name:   "filter",
			Usage:  "include or exclude files according to filtering rules",
			EnvVar: "PLUGIN_FILTER,FILTER",
		},
		cli.StringSliceFlag{
			Name:   "script",
			Usage:  "execute commands on the remote host after files are copied",
			EnvVar: "PLUGIN_SCRIPT,SCRIPT",
		},
		cli.StringFlag{
			Name:   "ssh-key",
			Usage:  "private ssh key",
			EnvVar: "PLUGIN_SSH_KEY,PLUGIN_KEY,SSH_KEY",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {
	plugin := Plugin{
		Config: Config{
			Hosts:     c.StringSlice("host"),
			User:      c.String("user"),
			Port:      c.Int("port"),
			Source:    c.String("source"),
			Target:    c.String("target"),
			Delete:    c.Bool("delete"),
			Recursive: c.Bool("recursive"),
			Include:   c.StringSlice("include"),
			Exclude:   c.StringSlice("exclude"),
			Filter:    c.StringSlice("filter"),
			Commands:  c.StringSlice("script"),
			Key:       c.String("ssh-key"),
		},
	}

	if len(plugin.Config.Key) == 0 {
		fmt.Println("No SSH key")
		os.Exit(1)
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
