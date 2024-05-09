package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "github.com/beevik/etree"
    "os"
)


func main() {
    if len(os.Args) < 2 {
        panic("Must pass shellid and commandid(optional) as parameter")
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

    Receive := etree.NewElement("rsp:Receive")
    DesiredStream := Receive.CreateElement("rsp:DesiredStream")
    if len(os.Args) == 3 {
        DesiredStream.CreateAttr("CommandId", os.Args[2])
    }
    DesiredStream.CreateText("stdout stderr")

    err, response := prot.Receive(resourceURI, Receive, &selectorset, nil)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    fmt.Println(response)
    fmt.Printf("\nDone\n")
}