package reader

import (
	"encoding/xml"
	"os"
)

type VisitorOptions struct {
	StreamFilename string
	IndexFilename  string
	Filter         *func(string) bool
}

func Visit(options *VisitorOptions, callback func(page Page) error) (err error) {
	file, err := os.Open(options.StreamFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = readIndex(options, func(indexes []*Index) error {
		matched := false
		if options.Filter != nil {
			for _, index := range indexes {
				if (*options.Filter)(index.Title) {
					matched = true
					break
				}
			}
		} else {
			matched = true
		}

		if !matched {
			return nil
		}

		innerXml, err := readStream(file, indexes[0])
		if err != nil {
			return err
		}

		xmlText := "<pages>" + innerXml + "</pages>"

		var ps *pages
		err = xml.Unmarshal([]byte(xmlText), &ps)
		if err != nil {
			return err
		}

		for _, page := range ps.Pages {
			if options.Filter != nil {
				if (*options.Filter)(page.Title) {
					err = callback(*page)
					if err != nil {
						return err
					}
				}
			} else {
				err = callback(*page)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}
