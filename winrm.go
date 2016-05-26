/*
Copyright 2013 Brice Figureau

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
	"os"

	"github.com/masterzen/winrm-cli/version"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "winrm"
	app.Version = version.Version
	app.Author = "WinRM contributors"
	app.Email = "https://github.com/masterzen/winrm-cli"
	app.Usage = "Command line tool to remotely execute commands on Windows machines"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
