package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
    "strings"
)

func main() {
    if len(os.Args) < 4 {
        panic("Must pass shellid, commandid and stream name as parameter")
    }
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()

    /* Test WMI Class */
    /*
    resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/Win32_OperatingSystem"
    */

    /* Test Schema enum with SelectorSet */
    /*
    resourceURI := "http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/*"
    selectorset := map[string]string {
        "__cimnamespace": cimNamespace,
        "ClassName": "",
    }
    */

    /* Test Shell */
    resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell"

    err, response_doc := prot.Enumerate(resourceURI, nil, nil, nil)
    if err != nil {
        fmt.Println(response_doc)
        panic(err)
    }

    fmt.Println(response_doc)
    fmt.Printf("\nDone\n")
}
