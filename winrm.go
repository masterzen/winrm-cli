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
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/masterzen/winrm"
	"github.com/mattn/go-isatty"
)

func main() {
	var (
		hostname  string
		user      string
		pass      string
		ntlm      bool
		kerberos  bool
		realm     string
		ccache    string
		cmd       string
		port      int
		encoded   bool
		https     bool
		insecure  bool
		cacert    string
		gencert   bool
		certsize  string
		timeout   string
	)

	flag.StringVar(&hostname, "hostname", "localhost", "winrm host")
	flag.StringVar(&user, "username", "vagrant", "winrm admin username")
	flag.StringVar(&pass, "password", "vagrant", "winrm admin password")
	flag.BoolVar(&ntlm, "ntlm", false, "use use ntlm auth")
	flag.BoolVar(&kerberos, "kerberos", false, "use kerberos auth")
	flag.StringVar(&realm, "realm", "", "kerberos realm")
	flag.StringVar(&ccache, "ccache", "", "kerberos credential cache file, usually /tmp/krb5cc_$UID")
	flag.BoolVar(&encoded, "encoded", false, "use base64 encoded password")
	flag.IntVar(&port, "port", 5985, "winrm port")
	flag.BoolVar(&https, "https", false, "use https")
	flag.BoolVar(&insecure, "insecure", false, "skip SSL validation")
	flag.StringVar(&cacert, "cacert", "", "CA certificate to use")
	flag.BoolVar(&gencert, "gencert", false, "Generate x509 client certificate to use with secure connections")
	flag.StringVar(&certsize, "certsize", "", "Priv RSA key between 512, 1024, 2048, 4096. Default :2048")
	flag.StringVar(&timeout, "timeout", "0s", "connection timeout")

	flag.Parse()

	if encoded {
		data, err := base64.StdEncoding.DecodeString(pass)
		check(err)
		pass = strings.TrimRight(string(data), "\r\n")
	}

	if gencert {
		cersize := pickSizeCert(certsize)
		config := CertConfig{
			Subject: pkix.Name{
				CommonName: "winrm client cert",
			},
			ValidFrom: time.Now(),
			ValidFor:  365 * 24 * time.Hour,
			SizeT:     cersize,
			Method:    RSA,
		}

		certPem, privPem, err := NewCert(config)
		check(err)
		err = ioutil.WriteFile("cert.cer", []byte(certPem), 0644)
		check(err)
		err = ioutil.WriteFile("priv.pem", []byte(privPem), 0644)
		check(err)
	} else {

		var (
			certBytes      []byte
			err            error
			connectTimeout time.Duration
		)

		if cacert != "" {
			certBytes, err = ioutil.ReadFile(cacert)
			check(err)
		} else {
			certBytes = nil
		}

		cmd = flag.Arg(0)

		if cmd == "" {
			check(errors.New("ERROR: Please enter the command to execute on the command line"))
		}

		connectTimeout, err = time.ParseDuration(timeout)
		check(err)

		endpoint := winrm.NewEndpoint(hostname, port, https, insecure, nil, certBytes, nil, connectTimeout)

		params := winrm.DefaultParameters

		if ntlm {
			params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }
		}

		if kerberos {
			proto := "http"
			if https == true {
				proto = "https"
			}

			params.TransportDecorator = func() winrm.Transporter { 
				return &winrm.ClientKerberos{
					Username:  user,
					Password:  pass,
					Hostname:  hostname,
					Realm:     realm,
					Port:      port,
					Proto:     proto,
					KrbConf:   "/etc/krb5.conf",
					KrbCCache: ccache,
					SPN:       fmt.Sprintf("HTTP/%s", hostname),
				}
			}
		}

		client, err := winrm.NewClientWithParameters(endpoint, user, pass, params)
		check(err)

		exitCode := 0
		if isatty.IsTerminal(os.Stdin.Fd()) {
			exitCode, err = client.Run(cmd, os.Stdout, os.Stderr)
		} else {
			exitCode, err = client.RunWithInput(cmd, os.Stdout, os.Stderr, os.Stdin)
		}
		check(err)

		os.Exit(exitCode)
	}
}

func pickSizeCert(size string) int {
	switch size {
	case "512":
		return 512
	case "1024":
		return 1024
	case "2048":
		return 2048
	case "4096":
		return 4096
	default:
		return 2048
	}
}

// generic check error func
func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
