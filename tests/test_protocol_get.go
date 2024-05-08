package main

import (
    "fmt"
    "gosamm/sammwr/protocol"
    "os"
)


func main() {
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
    fmt.Printf("Init Complete\n")


    /* Test WMI Class */
    /*
    resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/Win32_OperatingSystem"
    */

    /* Test Schema with SelectorSet */
    resourceURI := "http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/*"
    selectorset := map[string]string {
        "__cimnamespace": "root/cimv2",
        "ClassName": "Win32_OperatingSystem",
    }

    err, response_doc := prot.Get(resourceURI, &selectorset, nil)
    if err != nil {
        fmt.Println(response_doc)
        panic(err)
    }

    fmt.Println(response_doc)
    fmt.Printf("\nDone\n")
}