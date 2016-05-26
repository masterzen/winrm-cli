package command

import (
	"fmt"
	"runtime"

	"github.com/masterzen/winrm"
	"github.com/masterzen/winrm-cli/version"
	"github.com/urfave/cli"
)

func CmdVersion(c *cli.Context) error {
	fmt.Printf("winrm-cli version: %s\n", version.GetFullVersion())
	fmt.Printf("winrm library version: %s\n", winrm.GetFullVersion()) // not working yet
	fmt.Printf("Target OS/Arch: %s %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Built with Go Version: %s\n", runtime.Version())
	return nil
}
