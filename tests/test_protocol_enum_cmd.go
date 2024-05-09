package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
)

func main() {
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")
    resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell"

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    fmt.Printf("Init Complete\n")


    var ss *map[string]string

    if len(os.Args) > 1 {
        selectorset := map[string]string {
            "ShellId": os.Args[1],
        }
        ss = &selectorset
    } else {
        ss = nil
    }

    err, response_doc := prot.Enumerate(resourceURI, nil, ss, nil)
    if err != nil {
        fmt.Println(response_doc)
        panic(err)
    }

    fmt.Println(response_doc)
    fmt.Printf("\nDone\n")
}
