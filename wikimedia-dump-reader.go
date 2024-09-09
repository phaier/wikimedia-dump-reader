package main

import (
	"errors"
	"github.com/phaier/wikimedia-dump-reader/reader"
)

func main() {
	options := reader.VisitorOptions{
		StreamFilename: "data/jawiki-20240320-pages-articles-multistream.xml.bz2",
		IndexFilename:  "data/jawiki-20240320-pages-articles-multistream-index.txt.bz2",
		Filter:         nil,
	}

	err := reader.Visit(&options, func(page reader.Page) error {
		println(page.Title)

		return errors.New("stop")
	})

	if err != nil {
		panic(err)
	}
}
