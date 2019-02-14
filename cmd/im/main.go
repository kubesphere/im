/*
Copyright 2019 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
	"openpitrix.io/logger"

	"kubesphere.io/im/pkg/config"
	"kubesphere.io/im/pkg/service/im"
	"kubesphere.io/im/pkg/version"
)

var (
	appConfig *config.Config = nil
)

func main() {
	app := cli.NewApp()
	app.Name = "im"
	app.Usage = "provide im service."
	app.Version = version.GetVersionString()

	app.UsageText = `im [global options] command [options] [args...]

EXAMPLE:
   im gen-config
   im serve`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "config-im.json",
			Usage:  "im config file",
			EnvVar: "IM_CONFIG",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "im-service",
			EnvVar: "IM_HOST",
		},
	}

	app.Before = func(c *cli.Context) error {
		cfgPath := c.GlobalString("config")
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			appConfig = config.Default()
			ioutil.WriteFile(cfgPath, []byte(appConfig.JSONString()), 0666)
		} else {
			appConfig = config.MustLoad(c.GlobalString("config"))
		}

		logger.SetLevelByString(appConfig.LogLevel)
		return nil
	}

	app.Action = func(c *cli.Context) {
		serve(c)
	}

	app.Commands = []cli.Command{
		{
			Name:  "gen-config",
			Usage: "gen config file",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "json format (default)",
				},
				cli.BoolFlag{
					Name:  "toml",
					Usage: "toml format",
				},
				cli.BoolFlag{
					Name:  "yaml",
					Usage: "yaml format",
				},
			},

			Action: func(c *cli.Context) {
				switch {
				case c.Bool("json"):
					fmt.Println(config.Default().JSONString())
				case c.Bool("toml"):
					fmt.Println(config.Default().TOMLString())
				case c.Bool("yaml"):
					fmt.Println(config.Default().YAMLString())
				default:
					fmt.Println(config.Default().JSONString())
				}
				return
			},
		},

		{
			Name:  "serve",
			Usage: "run as service",
			Action: func(c *cli.Context) {
				serve(c)
			},
		},
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Fprintf(ctx.App.Writer, "not found '%v'!\n", command)
	}

	app.Run(os.Args)
}

func serve(c *cli.Context) {
	im.Serve(appConfig)
}
