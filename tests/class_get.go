package main

import (
	"github.com/samanamonitor/gosammwr/cim"
	"fmt"
	"os"
)

func main() {
    endpoint := os.Getenv("WR_ENDPOINT")
    username := os.Getenv("WR_USERNAME")
    password := os.Getenv("WR_PASSWORD")
    keytab_file := os.Getenv("WR_KEYTAB")
    if len(os.Args) < 3 {
        panic("Must pass namespace and classname as parameter and the command and arguments")
    }

    c, err := cim.NewClass(endpoint, username, password, keytab_file)
    if err != nil {
    	panic(err)
    }
    err = c.Get(os.Args[1], os.Args[2])
    if err != nil {
    	panic(err)
    }
    fmt.Println(c.ClassXML)
}