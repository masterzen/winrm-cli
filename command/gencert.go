package command

import (
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/masterzen/winrm-cli/certificate"
	"github.com/urfave/cli"
)

func CmdGencert(c *cli.Context) error {
	certsize := pickSizeCert(c.Int("certsize"))
	config := certificate.CertConfig{
		Subject: pkix.Name{
			CommonName: "winrm client cert",
		},
		ValidFrom: time.Now(),
		ValidFor:  365 * 24 * time.Hour,
		SizeT:     certsize,
		Method:    certificate.RSA,
	}

	certPem, privPem, err := certificate.NewCert(config)
	if err != nil {
		return createExitError("Generation failed: %s", err)
	}
	err = ioutil.WriteFile("cert.cer", []byte(certPem), 0644)
	if err != nil {
		return createExitError("Writing certificate failed: %s", err)
	}
	err = ioutil.WriteFile("priv.pem", []byte(privPem), 0644)
	if err != nil {
		return createExitError("Writing private key failed: %s", err)
	}

	return nil
}

func pickSizeCert(size int) int {
	switch size {
	case 512:
		return 512
	case 1024:
		return 1024
	case 2048:
		return 2048
	case 4096:
		return 4096
	default:
		return 2048
	}
}

func createExitError(msg string, err error) error {
	if err != nil {
		return cli.NewExitError(fmt.Sprintf(msg, err.Error()), 1)
	}
	return nil
}
