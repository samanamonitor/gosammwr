# Project Description
This project is intended to create a generic kerberos implementation for GOLANG. It is based on the source code in python named pykerberos. It uses native krb5 libraries instead of re-creating all kerberos implementation to reduce its size and get better compatibility.

Packages:
- libkrb5-dev
- krb5-user

# Development
`PROJECT_PATH=<path to the project>
ln -s ${PROJECT_PATH} /usr/local/go/src/gosamm`