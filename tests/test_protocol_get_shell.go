package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/sammwr/protocol"
    "os"
)


func main() {
    if len(os.Args) < 2 {
        panic("Must pass shellid as parameter")
    }
    endpoint := "http://smnnovmsfs01.samana.local:5985/wsman"
    keytab_file := "samanasvc2.keytab"
    var resourceURI string

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, nil, nil, &keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    fmt.Printf("Init Complete\n")


    resourceURI = "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"
    selectorset := map[string]string {
        "ShellId": os.Args[1],
    }

    err, response_doc := prot.Get(resourceURI, &selectorset, nil)
    if err != nil {
        response_doc.WriteTo(os.Stdout)
        fmt.Println()
        panic(err)
    }

    response_doc.WriteTo(os.Stdout)
    fmt.Printf("\nDone\n")
}