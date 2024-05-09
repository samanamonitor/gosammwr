package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "github.com/beevik/etree"
    "os"
)


func main() {
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




    shell := etree.NewElement("rsp:Shell")
    shell.CreateElement("rsp:InputStreams").CreateText("stdin")
    shell.CreateElement("rsp:OutputStreams").CreateText("stdout stderr")

    WorkingDirectory := "c:\\users"
    shell.CreateElement("rsp:WorkingDirectory").CreateText(WorkingDirectory)
    // Lifetime := "" iso8601_duration
    // shell.CreateElement("rsp:Lifetime").CreateText(Lifetime)
    IdleTimeOut := "P0Y0M0DT0H5M00S"
    shell.CreateElement("rsp:IdleTimeOut").CreateText(IdleTimeOut)
    /*
    Environment := map[string]string
    variables := shell.CreateElement("rsp:Environment")
    for key, val: range Environment {
        var := variables.CreateElement("rsp:Variable")
        var.CreateAttr("Name", key)
        var.CreateText(val)
    }
    */

    optionset := map[string]protocol.Option{}
    optionset["WINRS_NOPROFILE"] = protocol.Option{
        Value: "FALSE",
        Type: "xs:boolean",
    }
    optionset["WINRS_CODEPAGE"] = protocol.Option{
        Value: "437",
        Type: "xs:unsignedInt",
    }

    err, response := prot.Create(resourceURI, shell, &optionset)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    fmt.Println(response)
    fmt.Printf("\nDone\n")
}