package main

import (
	"github.com/beevik/etree"
	"fmt"
)

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("docs/get.xml"); err != nil {
		panic(err)
	}
	root := doc.SelectElement("Envelope")
	fmt.Println("ROOT element:", root.Tag)
	for _, header := range root.SelectElements("Header") {
		fmt.Println("CHILD element:", header.Tag)	
	}
	for _, t := range doc.FindElements("//a:MessageID") {
		fmt.Println("MessageID:", t.Text())
	}
	for _, t := range doc.FindElements("//w:ResourceURI") {
		fmt.Println("ResourceURI:", t.Text())
	}

}