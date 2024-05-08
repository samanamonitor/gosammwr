package main

import (
	"bytes"
	"fmt"
)

func main () {
	b := bytes.NewBuffer(nil)
	b.Write([]byte("This is a test"))
	var r [100]byte
	l, err := b.Read(r)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("(%d) %v\n", l, r)
}