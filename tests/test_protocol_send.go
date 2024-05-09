package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "github.com/beevik/etree"
    "os"
    "encoding/base64"
)


func main() {
    if len(os.Args) < 3 {
        panic("Must pass shellid and commandid as parameter")
    }


    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")
    resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    fmt.Printf("Init Complete\n")

    selectorset := map[string]string {
        "ShellId": os.Args[1],
    }

    Send := etree.NewElement("rsp:Send")
    Stream := Send.CreateElement("rsp:Stream")
    Stream.CreateAttr("Name", "stdin")
    Stream.CreateAttr("CommandId", os.Args[2])
    Stream.CreateAttr("End", "false")
    Stream.CreateText(base64.StdEncoding.EncodeToString([]byte("exit\r\n")))

    err, response := prot.Send(resourceURI, Send, &selectorset, nil)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    fmt.Println(response)
    fmt.Printf("\nDone\n")
}