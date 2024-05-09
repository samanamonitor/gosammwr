# Project Description
This project is intended to create a generic kerberos implementation for GOLANG. It is based on the source code in python named pykerberos. It uses native krb5 libraries instead of re-creating all kerberos implementation to reduce its size and get better compatibility.

Packages:
- libkrb5-dev
- krb5-user

# Development
`docker build -t gosammdev .`
`docker run -v $(pwd):/usr/src -idt --name gosammdev --dns 192.168.0.100 --dns 192.168.0.102 gosammdev /bin/bash`