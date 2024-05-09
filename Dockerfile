FROM golang:latest
RUN <<EOF
apt update
apt upgrade -y
apt install -y libkrb5-dev krb5-user
mkdir -p /go/pkg/mod/github.com/samanamonitor
ln -s /usr/src /go/pkg/mod/github.com/samanamonitor/gosammwr@v0.0.0-20240508230707-484b206977eb
EOF