package main

import (
	"fmt"
    "bytes"
    "mime/multipart"
    "net/textproto"
    "io"
)

func main() {
    data := []byte("test1234")


    f := bytes.NewBuffer(nil)
    w := multipart.NewWriter(f)
    w.SetBoundary("Encrypted Boundary")
    h := make(textproto.MIMEHeader)
    h.Set("Content-Type", "application/HTTP-SPNEGO-session-encrypted")
    value := fmt.Sprint("type=application/soap+xml;charset=UTF-8;Length=", len(data))
    h.Set("OriginalContent", value)
    w.CreatePart(h)
    h = make(textproto.MIMEHeader)
    h.Set("Content-Type", "application/octet-stream")
    part_writer, _ := w.CreatePart(h)
    io.Copy(part_writer, bytes.NewBuffer(data))
    w.Close()

    fmt.Println(f)
}
