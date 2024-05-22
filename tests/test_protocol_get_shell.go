package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
)


func main() {
    if len(os.Args) < 2 {
        panic("Must pass shellid as parameter")
    }
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    var resourceURI string

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    fmt.Printf("Init Complete\n")


    resourceURI = "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"
    selectorset := map[string]string {
        "ShellId": os.Args[1],
    }

    response_doc, err := prot.Get(resourceURI, &selectorset, nil)
    if err != nil {
        fmt.Println(response_doc)
        panic(err)
    }

    fmt.Println(response_doc)
    fmt.Printf("\nDone\n")
}