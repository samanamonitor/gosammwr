package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/protocol"
    "os"
    "bufio"
    "strings"
)


func main() {
    if len(os.Args) < 2 {
        panic("Must pass uuid as parameter")
    }
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    var resourceURI string
    var ec string

    prot, err := protocol.NewProtocol(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer prot.Close()
    if os.Args[1] == "-" {
        reader := bufio.NewReader(os.Stdin)
        input, _ := reader.ReadString('\n')
        temp := strings.Split(input, " ")
        if len(temp) < 2 {
            panic("Invalid number of inputs. <resourceuri> <enumerationcontext> expected.")
        }
        resourceURI = temp[0]
        ec = temp[1]
    } else {
        resourceURI = os.Args[1]
        ec = os.Args[2]
    }

    var response string
    for ; ec != ""; {
        ec, response, err = prot.Pull(resourceURI, ec, 5, nil, nil)
        if err != nil {
            fmt.Println(response)
            panic(err)
        }
    }

    fmt.Println(response)
}