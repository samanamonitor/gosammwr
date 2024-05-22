package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
    "strings"
)

/*
Examples:

bin/test_protocol_get http://schemas.dmtf.org/wbem/cim-xml/2/cim-schema/2/* __cimnamespace=root/cimv2 ClassName=Win32_OperatingSystem

bin/test_protocol_get http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/win32_diskdrive "DeviceId=\\\\.\\PHYSICALDRIVE2"

bin/test_protocol_get http://schemas.microsoft.com/wbem/wsman/1/windows/shell <shellid>
*/


func main() {
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")
    if len(os.Args) < 3 {
        panic("Must pass namespace and classname as parameter and the command and arguments")
    }

    var resourceURI string
    var selectorset *map[string]string

    resourceURI = os.Args[1]
    cmd_selectors := os.Args[2:]
    t_selectorset := map[string]string{}
    for i := 0; i < len(cmd_selectors); i += 1 {
        temp := strings.Split(cmd_selectors[i], "=")
        if len(temp) != 2 {
            panic("Invalid parameters. Selectors must be of type key=value")
        }
        t_selectorset[temp[0]] = temp[1]
    }
    selectorset = &t_selectorset

    prot := protocol.Protocol{}
    err := prot.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()

    response_doc, err := prot.Get(resourceURI, selectorset, nil)
    if err != nil {
        fmt.Println(response_doc)
        panic(err)
    }

    fmt.Println(response_doc)
}