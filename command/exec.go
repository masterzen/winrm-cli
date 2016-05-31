package command

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/masterzen/winrm"
	"github.com/urfave/cli"
)

func CmdExec(c *cli.Context) error {

	var (
		certBytes      []byte
		err            error
		connectTimeout time.Duration
	)

	if c.IsSet("cacert") {
		certBytes, err = ioutil.ReadFile(c.String("cacert"))
		if err != nil {
			return exitError(err)
		}
	}

	cmd := c.Args()[0]

	connectTimeout, err = time.ParseDuration(c.String("timeout"))
	if err != nil {
		return exitError(err)
	}

	endpoint := winrm.NewEndpointWithTimeout(c.String("hostname"), c.Int("port"), c.Bool("https"),
		c.Bool("insecure"), &certBytes, connectTimeout)

	client, err := winrm.NewClient(endpoint, c.String("username"), c.String("password"))
	if err != nil {
		return exitError(err)
	}

	exitCode, err := client.RunWithInput(cmd, os.Stdout, os.Stderr, os.Stdin)
	if err != nil {
		return exitError(err)
	}

	return cli.NewExitError("", exitCode)
}

func exitError(err error) error {
	return cli.NewExitError(err.Error(), 1)
}
