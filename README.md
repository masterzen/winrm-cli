# winrm-cli

This is a Go command-line executable to execute remote commands on Windows machines through
the use of WinRM/WinRS.

This tool support domain users via kerberos switch. You can either use a login/password combination, or an existing credential cache (obtained via "kinit")

## Contact

- Bugs: https://github.com/masterzen/winrm/issues


## Getting Started
WinRM is available on Windows Server 2008 and up. This project natively supports basic authentication for local accounts, see the steps in the next section on how to prepare the remote Windows machine for this scenario. The authentication model is pluggable, see below for an example on using Negotiate/NTLM authentication (e.g. for connecting to vanilla Azure VMs) or kerberos (for domain accounts)

### Preparing the remote Windows machine for Basic authentication
If you can not use kerberos, use basic authentication for local accounts. The remote windows system must be prepared for winrm:

_For a PowerShell script to do what is described below in one go, check [Richard Downer's blog](http://www.frontiertown.co.uk/2011/12/overthere-control-windows-from-java/)_

On the remote host, a PowerShell prompt, using the __Run as Administrator__ option and paste in the following lines:

    winrm quickconfig
    y
    winrm set winrm/config/service/Auth '@{Basic="true"}'
    winrm set winrm/config/service '@{AllowUnencrypted="true"}'
    winrm set winrm/config/winrs '@{MaxMemoryPerShellMB="1024"}'

__N.B.:__ The Windows Firewall needs to be running to run this command. See [Microsoft Knowledge Base article #2004640](http://support.microsoft.com/kb/2004640).

__N.B.:__ Do not disable Negotiate authentication as the windows `winrm` command itself uses this for internal authentication, and you risk getting a system where `winrm` doesn't work anymore.

__N.B.:__ The `MaxMemoryPerShellMB` option has no effects on some Windows 2008R2 systems because of a WinRM bug. Make sure to install the hotfix described [Microsoft Knowledge Base article #2842230](http://support.microsoft.com/kb/2842230) if you need to run commands that uses more than 150MB of memory.

For more information on WinRM, please refer to <a href="http://msdn.microsoft.com/en-us/library/windows/desktop/aa384426(v=vs.85).aspx">the online documentation at Microsoft's DevCenter</a>.

### Preparing the remote Windows machine for Kerberos authentication
The remote windows system must be prepared:

On the remote host, a PowerShell prompt, using the __Run as Administrator__ option and paste in the following lines:

    winrm quickconfig
    y
    winrm set winrm/config/service '@{AllowUnencrypted="true"}'
    winrm set winrm/config/winrs '@{MaxMemoryPerShellMB="1024"}'

### Building the winrm-cli executable

You can build winrm-cli from source:

```sh
git clone https://github.com/masterzen/winrm-cli
cd winrm-cli
make
```

This will generate a binary in the base directory called `./winrm`.

_Note_: you need go 1.5+. Please check your installation with

```
go version
```

## Command-line usage

Once built, you can run remote commands like this:

```sh
./winrm -hostname remote.domain.com -username "Administrator" -password "secret" "ipconfig /all"
```

### kerberos authentication

Either ensure you have a ticket granting ticket before running winrm (use kinit/klist), and run it with -ccache:

```sh
./winrm -kerberos -ccache /tmp/krb5cc_10007 -hostname remote.domain.com "ipconfig /all"
```

Or if you want to user login/password :

```sh
./winrm -kerberos -username test -password s3cr3t -hostname remote.domain.com -realm DOMAIN.COM "ipconfig /all"
```

## Docker

### Build image

```
docker build -t winrm .
```

### Usage

Once built, you can run remote commands like this:

```sh
docker run -it --rm winrm -hostname remote.domain.com -username "Administrator" -password "secret" "ipconfig /all"
```

