package protocol


import (
	"github.com/beevik/etree"
	"strings"
	"errors"
)

type Option struct {
	Type string
	Value string
}

type OptionSet map[string]Option

func (o *OptionSet) Add(n string, t string, v string) {
	(*o)[n] = Option{
		Type: t,
		Value: v,
	}
}

func (o *OptionSet) AddString(value string) error {
	temp := strings.Split(value, ",")
	if len(temp) != 3 {
		return errors.New("Selectors can only be set using format key=value")
	}
	(*o)[temp[0]] = Option{
		Type: temp[1],
		Value: temp[2],
	}
	return nil
}

func (o *OptionSet) ToElement(e *etree.Element) {
	if len(*o) == 0 {
		return
	}
	oset := e.CreateElement("w:OptionSet")
	for key, opt := range *o {
		option := oset.CreateElement("w:Option")
		option.CreateAttr("Name", key)
		option.CreateAttr("Type", opt.Type)
		option.CreateText(opt.Value)
	}
}
