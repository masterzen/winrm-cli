# winrm-cli

This is a Go command-line executable to execute remote commands on Windows machines through
the use of WinRM/WinRS.

This tool support domain users via kerberos switch. For this to work, you will need a valid TGT (use kinit), a valid krb5.conf and your target must have a resolvable reverse entry.

## Contact

- Bugs: https://github.com/masterzen/winrm/issues


## Getting Started
WinRM is available on Windows Server 2008 and up. This project natively supports basic authentication for local accounts, see the steps in the next section on how to prepare the remote Windows machine for this scenario. The authentication model is pluggable, see below for an example on using Negotiate/NTLM authentication (e.g. for connecting to vanilla Azure VMs).

### Preparing the remote Windows machine for Basic authentication (Not needed for Kerberos)
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

Ensure you have a ticket granting ticket before running winrm (use kinit/klist). Also make sure remote.domain.com can be resolved forward and reverse

```sh
./winrm -kerberos -hostname remote.domain.com "ipconfig /all"
```
These options can be set as environment variables before running winrm :  
KRB5_CONFIG points the kerberos config file. Default value is "/etc/krb5.conf"  
KRB5CCNAME points the credential cache file. Default value is "/tmp/krb5cc_${UID}"

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

