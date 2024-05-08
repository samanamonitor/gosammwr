package main

import (
    "encoding/xml"
    "fmt"
    "os"
    "strings"
    "reflect"
)

type Element struct {
	XMLName xml.Name
	Text *xml.CharData    `xml:",chardata"`
	Attr []xml.Attr       `xml:",any,attr"`
	Elements []Element    `xml:",any"`
	root *Element
	test string
	start *xml.StartElement
}

func (e *Element) String(level int) string {
	o := strings.Repeat(" ", level)
	o += fmt.Sprintf("Tag: %s\n", e.XMLName)
	o += strings.Repeat(" ", level)
	o += fmt.Sprintf("test: %s\n", e.test)
	o += strings.Repeat(" ", level)
	o += fmt.Sprintf("parent: %s\n", e.start)
	for i:=0; i < len(e.Attr); i++ {
		o += strings.Repeat(" ", level + 1)
		o += fmt.Sprintf("Attr - %s: %s\n", e.Attr[i].Name, e.Attr[i].Value)
	}
	if e.Text != nil {
		o += strings.Repeat(" ", level + 2)
		o += fmt.Sprintf("Text: %s\n", *e.Text)
	}
	for i:=0; i < len(e.Elements); i++ {
		o += e.Elements[i].String(level + 3)
	}
	return o
}

func (e *Element) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var t struct {
		XMLName xml.Name
		Text *xml.CharData    `xml:",chardata"`
		Attr []xml.Attr       `xml:",any,attr"`
		Elements []Element    `xml:",any"`
	}

	fmt.Printf("Type: %s\n", reflect.TypeOf(start))
	fmt.Printf("Token: %s\n", start)

	var cont bool
	for cont = true; cont ; {
		token, _ := d.Token()
		fmt.Printf("Type: %s\n", reflect.TypeOf(token))
		fmt.Printf("Token: %s\n", token)
		if token == nil {
			return nil
		}

	}
	if err := d.DecodeElement(&t, &start); err != nil {
		return err
	}
	fmt.Printf("name: %s\n", t.XMLName.Local)
	*e = Element{
		XMLName: t.XMLName,
		Attr: t.Attr,
		Text: t.Text,
		test: "this is a test",
		Elements: t.Elements,
		start: &start,
	}
	return nil
}

func printIndent(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
}

func printElement(e Element, indent int) {
	printIndent(indent)
	fmt.Printf("Tag: %s\n", e.XMLName)
	for i:=0; i < len(e.Attr); i++ {
		printIndent(indent + 1)
		fmt.Printf("Attr - %s: %s\n", e.Attr[i].Name, e.Attr[i].Value)
	}
	if e.Text != nil {
		printIndent(indent + 3)
		fmt.Printf("Text: %s\n", *e.Text)
	}
	for i:=0; i < len(e.Elements); i++ {
		printElement(e.Elements[i], indent+2)
	}

}

func main() {
	/*
	e := Envelope{
		S: "http://www.w3.org/2003/05/soap-envelope",
		A: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
		W: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
	}

	e.Header.Elements = append(e.Header.Elements, Element{
		XMLName: xml.Name{Local:"a:To"},
		Text: "http://windows-host:5985/wsman",
	})

	rt := Element{
		XMLName: xml.Name{Local: "a:ReplyTo"},
	}
	rt.Elements = append(rt.Elements, Element{
		XMLName: xml.Name{Local:"a:Address"},
		Text: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
	})
	e.Header.Elements = append(e.Header.Elements, rt)

	e.Header.Elements = append(e.Header.Elements, Element{
		XMLName: xml.Name{Local: "w:MaxEnvelopeSize"},
		Text: "512000",
	})

	e.Header.Elements = append(e.Header.Elements, Element{
		XMLName: xml.Name{Local:"a:MessageID"},
		Text: "uuid:FF2CF18E-0F75-4753-98F1-94F4D52D8DBF",
	})

	e.Header.Elements = append(e.Header.Elements, Element{
		XMLName: xml.Name{Local:"w:Locale"},
		Lang: "en-US",
		MustUnderstand: true,
	})

	e.Header.Elements = append(e.Header.Elements, Element{
		XMLName: xml.Name{Local:"p:DataLocale"},
		MustUnderstand: false,
		Lang: "en-US",
		Text: "",
	})

	e.Header.Elements = append(e.Header.Elements, Element{
		XMLName: xml.Name{Local: "w:ResourceURI"},
		Text: "http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/Win32_OperatingSystem",
	})

	e.Header.Elements = append(e.Header.Elements, Element{
		XMLName: xml.Name{Local: "a:Action"},
		Text: "http://schemas.xmlsoap.org/ws/2004/09/transfer/Get",
	})

	outtest, _ := xml.MarshalIndent(e, "", "  ")
    fmt.Println(xml.Header + string(outtest))

    fmt.Println()
    return
*/

/*
    coffee := &Plant{Id: 27, Name: "Coffee"}
    coffee.Origin = []string{"Ethiopia", "Brazil"}

    out, _ := xml.MarshalIndent(coffee, " ", "  ")
    fmt.Println(string(out))

    fmt.Println(xml.Header + string(out))

    var p Plant
    if err := xml.Unmarshal(out, &p); err != nil {
        panic(err)
    }
    fmt.Println(p)

    tomato := &Plant{Id: 81, Name: "Tomato"}
    tomato.Origin = []string{"Mexico", "California"}

    type Nesting struct {
        XMLName xml.Name `xml:"nesting"`
        Plants  []*Plant `xml:"parent>child>plant"`
    }

    nesting := &Nesting{}
    nesting.Plants = []*Plant{coffee, tomato}

    out, _ = xml.MarshalIndent(nesting, " ", "  ")
    fmt.Println(string(out))
*/
	message, err := os.ReadFile("docs/get.xml")
    if err != nil {
        panic(err)
    }
	var p Element
    err = xml.Unmarshal(message, &p)
    if err != nil {
    	panic(err)
    }
    /*
    printElement(p, 0)
    */
    return
    fmt.Print(p.String(0))
    fmt.Printf("\n\n")

    out, _ := xml.MarshalIndent(p, " ", "  ")
    fmt.Println(string(out))
}