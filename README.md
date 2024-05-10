# Project Description
This project is intended to create a generic kerberos implementation for GOLANG. It is based on the source code in python named pykerberos. It uses native krb5 libraries instead of re-creating all kerberos implementation to reduce its size and get better compatibility.

Packages:
- libkrb5-dev
- krb5-user

# Development
`docker build -t gosammdev .`
`DNS1=<change with kerberos dns server>
DNS2=<change with kerberos dns server>
docker run -v $(pwd):/usr/src -idt --name gosammdev --dns $DNS1 --dns $DNS2 gosammdev /bin/bash`

# Investigation tools
## Allow unencrypted messages on windows
```
winrm set winrm/config/client "@{AllowUnencrypted="true"}"
winrm set winrm/config/service "@{AllowUnencrypted="true"}"
```
## Send unencrypted messages with windows
```
$pssessionoption=New-PSSessionOption -NoEncryption
Enter-PSSession -ComputerName computer -SessionOption $pssessionoption
```
## Enable encrypted messages on Windows
```
winrm set winrm/config/client "@{AllowUnencrypted="false"}"
winrm set winrm/config/service "@{AllowUnencrypted="false"}"
```
# Shell
we need to create a connection for each stream and keep them alive. When a timeout is reached, the server will send a wsmanfault timeout. Then we need to re-send the receive command
we can send data over a separate connection

# Powershell
ResourceURI for powershell
http://schemas.microsoft.com/powershell/Microsoft.PowerShell
