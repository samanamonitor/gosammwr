package main

import (
	"github.com/beevik/etree"
	"os"
)

func main() {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	envelope := doc.CreateElement("s:Envelope")
	envelope.CreateAttr("xmlns:s", "http://www.w3.org/2003/05/soap-envelope")
	envelope.CreateAttr("xmlns:a", "http://schemas.xmlsoap.org/ws/2004/08/addressing")
	envelope.CreateAttr("xmlns:n", "http://schemas.xmlsoap.org/ws/2004/09/enumeration")
	envelope.CreateAttr("xmlns:w", "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd")
	envelope.CreateAttr("xmlns:p", "http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd")
	envelope.CreateAttr("xmlns:b", "http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd")

	header := envelope.CreateElement("s:Header")
	header.CreateElement("a:To").CreateText("http://windows-host:5985/wsman")
	address := header.CreateElement("a:ReplyTo").CreateElement("a:Address")
	address.CreateText("http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous")
	address.CreateAttr("s:mustUnderstand", "true")

	MaxEnvelopeSize := header.CreateElement("w:MaxEnvelopeSize")
	MaxEnvelopeSize.CreateAttr("s:mustUnderstand", "true")
	MaxEnvelopeSize.CreateText("512000")

	header.CreateElement("a:MessageID").CreateText("uuid:FF2CF18E-0F75-4753-98F1-94F4D52D8DBF")
	Locale := header.CreateElement("w:Locale") 
	Locale.CreateAttr("xml:lang", "en-US")
	Locale.CreateAttr("s:mustUnderstand", "false")

	DataLocale := header.CreateElement("p:DataLocale")
	DataLocale.CreateAttr("xml:lang", "en-US")
	DataLocale.CreateAttr("s:mustUnderstand", "false")

	ResourceURI := header.CreateElement("w:ResourceURI")
	ResourceURI.CreateAttr("s:mustUnderstand", "true")
	ResourceURI.CreateText("http://schemas.microsoft.com/wbem/wsman/1/wmi/root/cimv2/Win32_OperatingSystem")

	Action := header.CreateElement("a:Action")
	Action.CreateAttr("s:mustUnderstand", "true")
	Action.CreateText("http://schemas.xmlsoap.org/ws/2004/09/transfer/Get")

	envelope.CreateElement("s:Body")

	doc.WriteTo(os.Stdout)
}