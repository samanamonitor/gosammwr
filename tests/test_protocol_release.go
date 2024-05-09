package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
)


func main() {
    if len(os.Args) < 2 {
        panic("Must pass uuid as parameter")
    }
    endpoint := "http://smnnovmsfs01.samana.local:5985/wsman"
    username := ""
    password := ""
    keytab_file := "samanasvc2.keytab"
    baseURI := "http://schemas.microsoft.com/wbem/wsman/1/wmi"
    cimNamespace := "root/cimv2"
    className := "Win32_OperatingSystem"
    var resourceURI string

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    fmt.Printf("Init Complete\n")


    /* Test Resource Uri */
    resourceURI = baseURI + "/" + cimNamespace + "/" + className

    /* Test SelectorSet */
    /*
    resourceURI = "http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/*"
    selectorset := map[string]string {
        "__cimnamespace": cimNamespace,
        "ClassName": className,
    }
    */
    uuid := os.Args[1]

    err, response := prot.Release(resourceURI, uuid, nil)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    fmt.Println(response)
    fmt.Printf("\nDone\n")
}