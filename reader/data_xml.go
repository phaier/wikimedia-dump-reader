package reader

import "encoding/xml"

type Contributor struct {
	Username string `xml:"username"`
	Id       int    `xml:"id"`
}

type Revision struct {
	Id          int         `xml:"id"`
	ParentId    int         `xml:"parentid"`
	Timestamp   string      `xml:"timestamp"`
	Contributor Contributor `xml:"contributor"`
	Comment     string      `xml:"comment"`
	Model       string      `xml:"model"`
	Format      string      `xml:"format"`
	Text        string      `xml:"text"`
	Sha1        string      `xml:"sha1"`
}

type Page struct {
	Title    string   `xml:"title"`
	Id       int      `xml:"id"`
	Redirect string   `xml:"redirect"`
	Revision Revision `xml:"revision"`
}

type pages struct {
	XMLName xml.Name `xml:"pages"`
	Pages   []*Page  `xml:"page"`
}
