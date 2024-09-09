package reader

import (
	"bufio"
	"compress/bzip2"
	"os"
	"strconv"
	"strings"
)

type indexEntity struct {
	Offset int64
	Id     int64
	Title  string
}

func parseIndexEntity(line string) (*indexEntity, error) {
	index := strings.Index(line, ":")
	if index == -1 {
		panic("Invalid index line: " + line)
	}

	offset, err := strconv.ParseInt(line[:index], 10, 64)
	if err != nil {
		return nil, err
	}
	rest := line[index+1:]

	index = strings.Index(rest, ":")
	if index == -1 {
		panic("Invalid index line: " + line)
	}

	id, err := strconv.ParseInt(rest[:index], 10, 64)
	if err != nil {
		return nil, err
	}
	title := rest[index+1:]

	// ここで一行をパースする
	return &indexEntity{
		Offset: offset,
		Id:     id,
		Title:  title,
	}, nil
}

type Index struct {
	Offset int64
	Length int64
	Id     int64
	Title  string
}

func readIndex(options *VisitorOptions, callback func(indexes []*Index) error) error {
	file, err := os.Open(options.IndexFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bzip2.NewReader(file))

	var length = fileInfo.Size()
	var entities []*indexEntity

	for scanner.Scan() {
		line := scanner.Text()
		index, err := parseIndexEntity(line)
		if err != nil {
			return err
		}

		if len(entities) == 0 {
			entities = append(entities, index)
		} else {
			head := entities[0]
			if head.Offset != index.Offset {
				indexes := make([]*Index, len(entities))
				for i, entity := range entities {
					indexes[i] = &Index{
						Offset: entity.Offset,
						Length: index.Offset - entity.Offset,
						Id:     entity.Id,
						Title:  entity.Title,
					}
				}

				if err := callback(indexes); err != nil {
					return err
				}

				entities = nil
			} else {
				entities = append(entities, index)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	if len(entities) > 0 {
		indexes := make([]*Index, len(entities))
		for i, entity := range entities {
			indexes[i] = &Index{
				Offset: entity.Offset,
				Length: length - entity.Offset,
				Id:     entity.Id,
				Title:  entity.Title,
			}
		}

		if err := callback(indexes); err != nil {
			return err
		}
	}

	return nil
}
