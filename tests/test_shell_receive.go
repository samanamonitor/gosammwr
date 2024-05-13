package main

import (
    "fmt"
    "github.com/samanamonitor/gosammwr/shell"
    "os"
)


func main() {
    if len(os.Args) < 4 {
        panic("Must pass shellid, commandid and stream name as parameter")
    }
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")

    s := shell.Shell{}
    err := s.Init(endpoint, username, password, keytab_file)
    if err != nil {
        panic(err)
    }
    defer s.Cleanup()
    fmt.Printf("Init Complete\n")

    for ; ; {
        err, cs, response_streams := s.Receive(os.Args[1], os.Args[2], []string{os.Args[3]})
        if err != nil {
            panic(err)
        }
        fmt.Println(string(response_streams[0].Data))
        fmt.Print("CommandState: ")
        if cs.StateEnum == shell.CommandDone {
            fmt.Println("Done")
            fmt.Printf("ExitCode: %d\n", cs.ExitCode)
            break
        } else if cs.StateEnum == shell.CommandPending {
            fmt.Println("Pending")
        } else if cs.StateEnum == shell.CommandRunning {
            fmt.Println("Running")
        } else {
            panic("Invalid Command State")
        }

    }

    fmt.Printf("\nDone\n")
}