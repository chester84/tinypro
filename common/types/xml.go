// https://stackoverflow.com/questions/26371965/when-generating-an-xml-file-with-go-how-do-you-create-a-doctype-declaration
// https://stackoverflow.com/questions/14191596/how-to-create-a-cdata-node-of-xml-with-go
package types

import "encoding/xml"

type XmlCData string

func (n XmlCData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		S string `xml:",innerxml"`
	}{
		S: "<![CDATA[" + string(n) + "]]>",
	}, start)
}
