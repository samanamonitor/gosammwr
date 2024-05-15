package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "github.com/beevik/etree"
    "os"
    "bufio"
    "strings"
)

/*

example:

bin/test_protocol_create http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd

*/

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

    var resourceURI string

    if len(os.Args) != 2 {
        panic("Invalid number of parameters. Expecting resourceuri or -")
    }
    if os.Args[1] == "-" {
        reader := bufio.NewReader(os.Stdin)
        input, _ := reader.ReadString('\n')
        resourceURI = strings.Trim(input, "\n")

    } else {
        resourceURI = os.Args[1]
    }

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

    optionset := map[string]protocol.Option{
        "WINRS_NOPROFILE": { Value: "FALSE", Type: "xs:boolean" },
        "WINRS_CODEPAGE": { Value: "437", Type: "xs:unsignedInt" },
    }

    err, response := prot.Create(resourceURI, shell, &optionset)
    if err != nil {
        fmt.Println(response)
        panic(err)
    }

    Body := etree.NewDocument()
    Body.ReadFromString(response)
    ShellId := Body.FindElement("//Shell/ShellId")
    if ShellId == nil {
        fmt.Println(response)
        panic("Invalid response. Missing ShellId")
    }
    fmt.Println(resourceURI, ShellId.Text())
}