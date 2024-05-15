package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
    "github.com/beevik/etree"
)


func main() {
    if len(os.Args) < 3 {
        panic("Must pass shellid as parameter and the command and arguments")
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

    resourceURI := "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd"

    selectorset := map[string]string{
        "ShellId": os.Args[1],
    }
    optionset := map[string]protocol.Option {
        "WINRS_CONSOLEMODE_STDIN": { Value: "TRUE", Type: "xs:boolean" },
        "WINRS_SKIP_CMD_SHELL": { Value: "TRUE", Type: "xs:boolean" },
    }

    command := os.Args[2:]
    if len(command) < 1 {
        panic("command array must have at least 1 element")
    }
    CommandLine := etree.NewElement("rsp:CommandLine")
    CommandLine.CreateElement("rsp:Command").CreateText(command[0])
    for i := 1; i < len(command); i++ {
        CommandLine.CreateElement("rsp:Arguments").CreateText(command[i])
    }

    err, body_str := prot.Command(resourceURI, CommandLine, &selectorset, &optionset)
    if err != nil {
        panic(err)
    }
    Body := etree.NewDocument()
    Body.ReadFromString(body_str)
    CommandId := Body.FindElement("//CommandId")
    fmt.Println(CommandId.Text())
}
