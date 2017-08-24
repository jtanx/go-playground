package main

import (
	"encoding/xml"
	"fmt"
)

type A struct {
	Opt string
	Arg string
}

type B A

func (b *B) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}
	b.Opt = str
	b.Arg = "custom"
	return nil
}

func main() {
	var a A
	var b B
	err := xml.Unmarshal([]byte(`<AAA><Opt>where</Opt><Arg>are</Arg>aaa</AAA>`), &a)
	fmt.Printf("%v %v\n", a, err)
	err = xml.Unmarshal([]byte(`<AAA><Opt>where</Opt><Arg>are</Arg>aaa</AAA>`), &b)
	fmt.Printf("%v %v\n", b, err)
}
