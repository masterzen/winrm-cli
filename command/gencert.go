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
	var (
		method   certificate.KeyType
		cn       = c.String("cn")
		certsize int
		err      error
	)

	sz := c.Int("certsize")
	if c.Bool("rsa") {
		method = certificate.RSA
		certsize, err = pickRSACertSize(sz)
		if err != nil {
			return exitError(fmt.Errorf("Cert size: %s", err))
		}
	} else if c.Bool("ecdsa") {
		method = certificate.ECDSA
		certsize, err = pickECDSAcertSize(sz)
		if err != nil {
			return exitError(fmt.Errorf("Cert size: %s", err))
		}
	} else {
		return exitError(fmt.Errorf("Unsupported encryption: %s",
			"Please specify the encryption RSA or ECDSA"))
	}

	config := certificate.CertConfig{
		Subject: pkix.Name{
			CommonName: cn,
		},
		ValidFrom: time.Now(),
		ValidFor:  365 * 24 * time.Hour,
		SizeT:     certsize,
		Method:    method,
	}

	certPem, privPem, err := certificate.NewCert(config)
	if err != nil {
		return exitError(fmt.Errorf("Generation failed: %s", err))
	}

	err = ioutil.WriteFile("cert.cer", []byte(certPem), 0644)
	if err != nil {
		return exitError(fmt.Errorf("Writing certificate failed: %s", err))
	}
	err = ioutil.WriteFile("priv.pem", []byte(privPem), 0644)
	if err != nil {
		return exitError(fmt.Errorf("Writing private key failed: %s", err))
	}

	return nil
}

// test if it's a valid rsa size
func pickRSACertSize(size int) (int, error) {
	switch size {
	case 512:
		return size, nil
	case 1024:
		return size, nil
	case 2048:
		return size, nil
	case 4096:
		return size, nil
	default:
		return 0, fmt.Errorf("Unsupported RSA cert size")
	}
}

// test if it's a valid ecdsa size
func pickECDSAcertSize(size int) (int, error) {
	switch size {
	case 224:
		return size, nil
	case 256:
		return size, nil
	case 384:
		return size, nil
	case 521:
		return size, nil
	default:
		return 0, fmt.Errorf("Unsupported ECDSA cert size")
	}

}
