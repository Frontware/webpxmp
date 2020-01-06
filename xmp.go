package xmp

import "encoding/xml"

type Profile struct {
	Text      string  `xml:",chardata"`
	Name      string  `xml:"name"`
	Timestamp int64   `xml:"timestamp"`
	Location  string  `xml:"location"`
	Lat       float64 `xml:"lat"`
	Long      float64 `xml:"long"`
}

type Xmpmeta struct {
	XMLName xml.Name `xml:"xmpmeta"`
	Text    string   `xml:",chardata"`
	X       string   `xml:"x,attr"`
	Xmptk   string   `xml:"xmptk,attr"`
	RDF     struct {
		Text        string `xml:",chardata"`
		Rdf         string `xml:"rdf,attr"`
		Description struct {
			Text    string  `xml:",chardata"`
			About   string  `xml:"about,attr"`
			Weladee string  `xml:"weladee,attr"`
			Profile Profile `xml:"profile"`
		} `xml:"Description"`
	} `xml:"RDF"`
}
