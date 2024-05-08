package main

import (
	"github.com/beevik/etree"
	"fmt"
)

func main() {
	/*
	data := "<n:PullResponse><n:Items/><n:EndOfSequence/></n:PullResponse>"
	*/
	data := "<n:PullResponse><n:Items><rsp:Shell xmlns:rsp=\"http://schemas.microsoft.com/wbem/wsman/1/windows/shell\"><rsp:ShellId>242FA19C-9586-49D4-84CD-A28D1D3A402F</rsp:ShellId><rsp:ResourceUri>http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</rsp:ResourceUri><rsp:Owner>SAMANA\\samanasvc2</rsp:Owner><rsp:ClientIP>192.168.69.11</rsp:ClientIP><rsp:WorkingDirectory>c:\\users</rsp:WorkingDirectory><rsp:IdleTimeOut>PT7200.000S</rsp:IdleTimeOut><rsp:InputStreams>stdin</rsp:InputStreams><rsp:OutputStreams>stdout stderr</rsp:OutputStreams><rsp:ShellRunTime>P0DT0H40M30S</rsp:ShellRunTime><rsp:ShellInactivity>P0DT0H34M10S</rsp:ShellInactivity></rsp:Shell></n:Items><n:EndOfSequence/></n:PullResponse>"
	doc := etree.NewDocument()
	doc.ReadFromBytes([]byte(data))
	newdoc := etree.NewDocument()
	items := newdoc.CreateElement("Items")
	for _, e := range doc.FindElements("//PullResponse/Items/*") {
		items.AddChild(e)
	}
	str, _ := newdoc.WriteToString()
	fmt.Println(str)
}

